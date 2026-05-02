package notifier

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

const (
	acsAPIVersion = "2023-03-31"
	acsScope      = "https://communication.azure.com/.default"
)

// ACSConfig configures the Azure Communication Services Email notifier.
type ACSConfig struct {
	Endpoint         string // e.g. https://carrental-comm.communication.azure.com
	Sender           string // e.g. DoNotReply@<verified-domain>
	FromName         string
	AuthMode         string // "managed_identity" or "connection_string"
	ConnectionString string // only required when AuthMode == "connection_string"
}

// ACSNotifier sends transactional email through Azure Communication Services.
// Auth is either Managed Identity (production) or Connection String (dev).
type ACSNotifier struct {
	endpoint   string
	sender     string
	fromName   string
	httpClient *http.Client

	// One of these auth strategies is set, never both.
	credential azcore.TokenCredential
	accessKey  []byte // raw bytes for HMAC signing
}

func NewACS(cfg ACSConfig) (*ACSNotifier, error) {
	if cfg.Endpoint == "" {
		return nil, errors.New("acs: endpoint is required")
	}
	if cfg.Sender == "" {
		return nil, errors.New("acs: sender is required")
	}

	n := &ACSNotifier{
		endpoint:   strings.TrimRight(cfg.Endpoint, "/"),
		sender:     cfg.Sender,
		fromName:   cfg.FromName,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	switch strings.ToLower(cfg.AuthMode) {
	case "", "managed_identity":
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return nil, fmt.Errorf("acs: managed identity: %w", err)
		}
		n.credential = cred
	case "connection_string":
		endpoint, key, err := parseConnectionString(cfg.ConnectionString)
		if err != nil {
			return nil, err
		}
		// Connection-string endpoint overrides explicit endpoint when present
		// to avoid mismatches.
		if endpoint != "" {
			n.endpoint = strings.TrimRight(endpoint, "/")
		}
		n.accessKey = key
	default:
		return nil, fmt.Errorf("acs: unknown auth mode %q", cfg.AuthMode)
	}

	return n, nil
}

// Send delivers msg via the ACS Email API. A non-2xx response yields an
// error so the caller (RabbitMQ consumer) can NACK for redelivery.
func (n *ACSNotifier) Send(ctx context.Context, msg Message) error {
	body, err := json.Marshal(acsRequest{
		SenderAddress: n.sender,
		Content: acsContent{
			Subject:   msg.Subject,
			PlainText: msg.TextBody,
			HTML:      msg.HTMLBody,
		},
		Recipients: acsRecipients{
			To: []acsRecipient{{Address: msg.To, DisplayName: ""}},
		},
	})
	if err != nil {
		return fmt.Errorf("acs: marshal: %w", err)
	}

	reqURL := fmt.Sprintf("%s/emails:send?api-version=%s", n.endpoint, acsAPIVersion)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("acs: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("repeatability-request-id", msg.To+"-"+time.Now().UTC().Format(time.RFC3339Nano))
	req.Header.Set("repeatability-first-sent", time.Now().UTC().Format(http.TimeFormat))

	if err := n.authorize(ctx, req, body); err != nil {
		return err
	}

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("acs: http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		slog.Info("acs email accepted",
			"status", resp.StatusCode,
			"to", redactEmail(msg.To),
			"subject", msg.Subject,
			"operation_location", resp.Header.Get("Operation-Location"))
		return nil
	}

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return fmt.Errorf("acs: send failed: status=%d body=%s", resp.StatusCode, string(respBody))
}

// authorize attaches the appropriate auth header to req. Two modes:
//   - Managed Identity: Bearer token from azidentity.
//   - Connection String: HMAC-SHA256 over date+host+content-hash.
func (n *ACSNotifier) authorize(ctx context.Context, req *http.Request, body []byte) error {
	if n.credential != nil {
		tok, err := n.credential.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{acsScope}})
		if err != nil {
			return fmt.Errorf("acs: get token: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+tok.Token)
		return nil
	}

	// HMAC-SHA256 signing per ACS Connection-String auth scheme.
	now := time.Now().UTC().Format(http.TimeFormat)
	hash := sha256.Sum256(body)
	contentHash := base64.StdEncoding.EncodeToString(hash[:])

	parsed, err := url.Parse(req.URL.String())
	if err != nil {
		return fmt.Errorf("acs: parse url: %w", err)
	}
	pathAndQuery := parsed.RequestURI()
	host := parsed.Host

	stringToSign := fmt.Sprintf("%s\n%s\n%s;%s;%s", req.Method, pathAndQuery, now, host, contentHash)
	mac := hmac.New(sha256.New, n.accessKey)
	mac.Write([]byte(stringToSign))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	req.Header.Set("x-ms-date", now)
	req.Header.Set("x-ms-content-sha256", contentHash)
	req.Header.Set("Authorization",
		"HMAC-SHA256 SignedHeaders=x-ms-date;host;x-ms-content-sha256&Signature="+sig)
	return nil
}

// parseConnectionString accepts forms like
//
//	endpoint=https://x.communication.azure.com;accesskey=BASE64KEY
//
// and returns the endpoint and the raw decoded key bytes.
func parseConnectionString(s string) (endpoint string, key []byte, err error) {
	if s == "" {
		return "", nil, errors.New("acs: empty connection string")
	}
	for _, part := range strings.Split(s, ";") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(kv[0])) {
		case "endpoint":
			endpoint = strings.TrimSpace(kv[1])
		case "accesskey":
			raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(kv[1]))
			if err != nil {
				return "", nil, fmt.Errorf("acs: bad accesskey: %w", err)
			}
			key = raw
		}
	}
	if endpoint == "" || len(key) == 0 {
		return "", nil, errors.New("acs: connection string missing endpoint or accesskey")
	}
	return endpoint, key, nil
}

// redactEmail turns "john.doe@example.com" into "j***@example.com" for logs.
func redactEmail(addr string) string {
	at := strings.IndexByte(addr, '@')
	if at <= 1 {
		return "***"
	}
	return addr[:1] + "***" + addr[at:]
}

// ---- ACS REST payload types ------------------------------------------------

type acsRequest struct {
	SenderAddress string        `json:"senderAddress"`
	Content       acsContent    `json:"content"`
	Recipients    acsRecipients `json:"recipients"`
}

type acsContent struct {
	Subject   string `json:"subject"`
	PlainText string `json:"plainText"`
	HTML      string `json:"html"`
}

type acsRecipients struct {
	To []acsRecipient `json:"to"`
}

type acsRecipient struct {
	Address     string `json:"address"`
	DisplayName string `json:"displayName,omitempty"`
}

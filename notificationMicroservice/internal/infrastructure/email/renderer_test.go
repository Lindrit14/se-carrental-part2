package email

import (
	"strings"
	"testing"
)

func TestRendererAllTemplates(t *testing.T) {
	r, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer: %v", err)
	}

	cases := []struct {
		name string
		data map[string]any
		// must contain at least these substrings in the rendered text body
		expectInText []string
	}{
		{
			name: "welcome",
			data: map[string]any{"Email": "alice@example.com"},
			expectInText: []string{
				"Welcome to Car Rental",
				"alice@example.com",
			},
		},
		{
			name: "password_reset",
			data: map[string]any{
				"Email":     "alice@example.com",
				"ResetURL":  "https://app.example.com/password/confirm?token=xyz",
				"ExpiresAt": "2026-05-02T10:00:00Z",
			},
			expectInText: []string{
				"Password Reset",
				"https://app.example.com/password/confirm?token=xyz",
				"2026-05-02T10:00:00Z",
			},
		},
		{
			name: "booking_confirmation",
			data: map[string]any{
				"BookingID":     "BK-123",
				"CarID":         "CAR-7",
				"StartDate":     "2026-06-01",
				"EndDate":       "2026-06-05",
				"TotalAmount":   "199.50",
				"TotalCurrency": "EUR",
			},
			expectInText: []string{"BK-123", "CAR-7", "2026-06-01", "199.50", "EUR"},
		},
		{
			name: "booking_cancellation",
			data: map[string]any{"BookingID": "BK-123"},
			expectInText: []string{"BK-123", "Cancelled"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			html, text, err := r.Render(tc.name, tc.data)
			if err != nil {
				t.Fatalf("Render: %v", err)
			}
			if html == "" || text == "" {
				t.Fatalf("empty body: html=%d text=%d", len(html), len(text))
			}
			for _, want := range tc.expectInText {
				if !strings.Contains(text, want) {
					t.Errorf("text missing %q\n--- got:\n%s", want, text)
				}
			}
		})
	}
}

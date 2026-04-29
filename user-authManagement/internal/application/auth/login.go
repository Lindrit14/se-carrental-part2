package auth

import (
	"context"
	"errors"
	"time"

	"github.com/lindritprekaj/user-authmanagement/internal/application/ports"
	"github.com/lindritprekaj/user-authmanagement/internal/domain/token"
	domainuser "github.com/lindritprekaj/user-authmanagement/internal/domain/user"
)

type LoginInput struct {
	Email     string
	Password  string
	IP        string
	UserAgent string
}

type LoginOutput struct {
	AccessToken     string
	AccessExpiresAt time.Time
	RefreshToken    string
	RefreshExpiresAt time.Time
	UserID          string
}

type LoginUseCase struct {
	users     domainuser.Repository
	tokens    token.Repository
	hasher    ports.PasswordHasher
	tokenSvc  ports.TokenService
	publisher ports.EventPublisher
	clock     ports.Clock
	ids       ports.IDGenerator
	policy    Policy
}

func NewLoginUseCase(
	users domainuser.Repository,
	tokens token.Repository,
	hasher ports.PasswordHasher,
	tokenSvc ports.TokenService,
	publisher ports.EventPublisher,
	clock ports.Clock,
	ids ports.IDGenerator,
	policy Policy,
) *LoginUseCase {
	return &LoginUseCase{
		users: users, tokens: tokens, hasher: hasher,
		tokenSvc: tokenSvc, publisher: publisher,
		clock: clock, ids: ids, policy: policy,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, in LoginInput) (*LoginOutput, error) {
	u, err := uc.users.FindByEmail(ctx, domainuser.NormalizeEmail(in.Email))
	if err != nil {
		if errors.Is(err, domainuser.ErrNotFound) {
			return nil, domainuser.ErrInvalidCredentials
		}
		return nil, err
	}
	if err := uc.hasher.Verify(in.Password, u.PasswordHash); err != nil {
		return nil, domainuser.ErrInvalidCredentials
	}

	roles := make([]string, 0, len(u.Roles))
	for _, r := range u.Roles {
		roles = append(roles, string(r))
	}

	access, accessExp, err := uc.tokenSvc.IssueAccessToken(u.ID, roles)
	if err != nil {
		return nil, err
	}

	refreshPlain, err := uc.tokenSvc.NewRefreshToken()
	if err != nil {
		return nil, err
	}
	now := uc.clock.Now()
	rt := &token.RefreshToken{
		ID:        uc.ids.New(),
		UserID:    u.ID,
		TokenHash: uc.tokenSvc.HashRefreshToken(refreshPlain),
		IssuedAt:  now,
		ExpiresAt: now.Add(uc.policy.RefreshTTL),
	}
	if err := uc.tokens.Create(ctx, rt); err != nil {
		return nil, err
	}

	_ = uc.publisher.Publish(ctx, ports.Event{
		ID:         uc.ids.New(),
		Type:       EventUserLogin,
		Version:    EventVersion,
		OccurredAt: now.UTC().Format(rfc3339),
		RoutingKey: EventUserLogin,
		Data: UserLoginData{
			UserID:    u.ID,
			IP:        in.IP,
			UserAgent: in.UserAgent,
		},
	})

	return &LoginOutput{
		AccessToken:      access,
		AccessExpiresAt:  accessExp,
		RefreshToken:     refreshPlain,
		RefreshExpiresAt: rt.ExpiresAt,
		UserID:           u.ID,
	}, nil
}

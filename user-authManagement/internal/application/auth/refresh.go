package auth

import (
	"context"
	"errors"
	"time"

	"github.com/lindritprekaj/user-authmanagement/internal/application/ports"
	"github.com/lindritprekaj/user-authmanagement/internal/domain/token"
	domainuser "github.com/lindritprekaj/user-authmanagement/internal/domain/user"
)

type RefreshInput struct {
	RefreshToken string
}

type RefreshOutput struct {
	AccessToken      string
	AccessExpiresAt  time.Time
	RefreshToken     string
	RefreshExpiresAt time.Time
}

type RefreshUseCase struct {
	users    domainuser.Repository
	tokens   token.Repository
	tokenSvc ports.TokenService
	clock    ports.Clock
	ids      ports.IDGenerator
	policy   Policy
}

func NewRefreshUseCase(
	users domainuser.Repository,
	tokens token.Repository,
	tokenSvc ports.TokenService,
	clock ports.Clock,
	ids ports.IDGenerator,
	policy Policy,
) *RefreshUseCase {
	return &RefreshUseCase{
		users: users, tokens: tokens, tokenSvc: tokenSvc,
		clock: clock, ids: ids, policy: policy,
	}
}

func (uc *RefreshUseCase) Execute(ctx context.Context, in RefreshInput) (*RefreshOutput, error) {
	hash := uc.tokenSvc.HashRefreshToken(in.RefreshToken)
	current, err := uc.tokens.FindByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	now := uc.clock.Now()
	if !current.IsActive(now) {
		// Reuse of a revoked token is suspicious — revoke the entire family.
		if current.RevokedAt != nil {
			_ = uc.tokens.RevokeAllForUser(ctx, current.UserID)
		}
		return nil, token.ErrRevoked
	}

	u, err := uc.users.FindByID(ctx, current.UserID)
	if err != nil {
		if errors.Is(err, domainuser.ErrNotFound) {
			return nil, domainuser.ErrInvalidCredentials
		}
		return nil, err
	}

	roles := make([]string, 0, len(u.Roles))
	for _, r := range u.Roles {
		roles = append(roles, string(r))
	}
	access, accessExp, err := uc.tokenSvc.IssueAccessToken(u.ID, roles)
	if err != nil {
		return nil, err
	}

	newPlain, err := uc.tokenSvc.NewRefreshToken()
	if err != nil {
		return nil, err
	}
	newRT := &token.RefreshToken{
		ID:        uc.ids.New(),
		UserID:    u.ID,
		TokenHash: uc.tokenSvc.HashRefreshToken(newPlain),
		IssuedAt:  now,
		ExpiresAt: now.Add(uc.policy.RefreshTTL),
	}
	if err := uc.tokens.Create(ctx, newRT); err != nil {
		return nil, err
	}

	current.Revoke(now, newRT.ID)
	if err := uc.tokens.Update(ctx, current); err != nil {
		return nil, err
	}

	return &RefreshOutput{
		AccessToken:      access,
		AccessExpiresAt:  accessExp,
		RefreshToken:     newPlain,
		RefreshExpiresAt: newRT.ExpiresAt,
	}, nil
}

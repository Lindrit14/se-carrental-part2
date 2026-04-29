package auth

import (
	"context"
	"errors"

	"github.com/lindritprekaj/user-authmanagement/internal/application/ports"
	"github.com/lindritprekaj/user-authmanagement/internal/domain/token"
	domainuser "github.com/lindritprekaj/user-authmanagement/internal/domain/user"
)

type RequestPasswordResetInput struct {
	Email string
}

type RequestPasswordResetUseCase struct {
	users     domainuser.Repository
	resets    token.PasswordResetRepository
	tokenSvc  ports.TokenService
	publisher ports.EventPublisher
	clock     ports.Clock
	ids       ports.IDGenerator
	policy    Policy
}

func NewRequestPasswordResetUseCase(
	users domainuser.Repository,
	resets token.PasswordResetRepository,
	tokenSvc ports.TokenService,
	publisher ports.EventPublisher,
	clock ports.Clock,
	ids ports.IDGenerator,
	policy Policy,
) *RequestPasswordResetUseCase {
	return &RequestPasswordResetUseCase{
		users: users, resets: resets, tokenSvc: tokenSvc,
		publisher: publisher, clock: clock, ids: ids, policy: policy,
	}
}

// Execute always returns nil error to avoid leaking which emails are registered.
// On unknown email it silently succeeds; only logs/metrics distinguish.
func (uc *RequestPasswordResetUseCase) Execute(ctx context.Context, in RequestPasswordResetInput) error {
	email := domainuser.NormalizeEmail(in.Email)
	u, err := uc.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domainuser.ErrNotFound) {
			return nil
		}
		return err
	}

	plain, err := uc.tokenSvc.NewRefreshToken()
	if err != nil {
		return err
	}
	now := uc.clock.Now()
	rec := &token.PasswordReset{
		ID:        uc.ids.New(),
		UserID:    u.ID,
		TokenHash: uc.tokenSvc.HashRefreshToken(plain),
		ExpiresAt: now.Add(uc.policy.PasswordResetTTL),
	}
	if err := uc.resets.Create(ctx, rec); err != nil {
		return err
	}

	_ = uc.publisher.Publish(ctx, ports.Event{
		ID:         uc.ids.New(),
		Type:       EventUserPasswordResetRequested,
		Version:    EventVersion,
		OccurredAt: now.UTC().Format(rfc3339),
		RoutingKey: EventUserPasswordResetRequested,
		Data: UserPasswordResetRequestedData{
			UserID:     u.ID,
			Email:      u.Email,
			ResetToken: plain,
			ExpiresAt:  rec.ExpiresAt.UTC().Format(rfc3339),
		},
	})
	return nil
}

type ConfirmPasswordResetInput struct {
	ResetToken  string
	NewPassword string
}

type ConfirmPasswordResetUseCase struct {
	users    domainuser.Repository
	resets   token.PasswordResetRepository
	tokens   token.Repository
	hasher   ports.PasswordHasher
	tokenSvc ports.TokenService
	clock    ports.Clock
	policy   Policy
}

func NewConfirmPasswordResetUseCase(
	users domainuser.Repository,
	resets token.PasswordResetRepository,
	tokens token.Repository,
	hasher ports.PasswordHasher,
	tokenSvc ports.TokenService,
	clock ports.Clock,
	policy Policy,
) *ConfirmPasswordResetUseCase {
	return &ConfirmPasswordResetUseCase{
		users: users, resets: resets, tokens: tokens,
		hasher: hasher, tokenSvc: tokenSvc, clock: clock, policy: policy,
	}
}

func (uc *ConfirmPasswordResetUseCase) Execute(ctx context.Context, in ConfirmPasswordResetInput) error {
	if err := uc.policy.ValidatePassword(in.NewPassword); err != nil {
		return err
	}
	hash := uc.tokenSvc.HashRefreshToken(in.ResetToken)
	rec, err := uc.resets.FindByHash(ctx, hash)
	if err != nil {
		return err
	}
	now := uc.clock.Now()
	if !rec.IsActive(now) {
		return token.ErrExpired
	}

	u, err := uc.users.FindByID(ctx, rec.UserID)
	if err != nil {
		return err
	}
	newHash, err := uc.hasher.Hash(in.NewPassword)
	if err != nil {
		return err
	}
	u.ChangePassword(newHash, now)
	if err := uc.users.Update(ctx, u); err != nil {
		return err
	}
	if err := uc.resets.MarkUsed(ctx, rec.ID, now); err != nil {
		return err
	}
	// Invalidate all existing sessions on password change.
	return uc.tokens.RevokeAllForUser(ctx, u.ID)
}

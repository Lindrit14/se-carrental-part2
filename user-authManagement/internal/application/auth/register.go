// Package auth contains the authentication use cases.
package auth

import (
	"context"
	"errors"
	"net/mail"

	"github.com/lindritprekaj/user-authmanagement/internal/application/ports"
	domainuser "github.com/lindritprekaj/user-authmanagement/internal/domain/user"
)

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterOutput struct {
	UserID string
	Email  string
}

type RegisterUseCase struct {
	users     domainuser.Repository
	hasher    ports.PasswordHasher
	publisher ports.EventPublisher
	clock     ports.Clock
	ids       ports.IDGenerator
	policy    Policy
}

func NewRegisterUseCase(
	users domainuser.Repository,
	hasher ports.PasswordHasher,
	publisher ports.EventPublisher,
	clock ports.Clock,
	ids ports.IDGenerator,
	policy Policy,
) *RegisterUseCase {
	return &RegisterUseCase{
		users:     users,
		hasher:    hasher,
		publisher: publisher,
		clock:     clock,
		ids:       ids,
		policy:    policy,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, in RegisterInput) (*RegisterOutput, error) {
	email := domainuser.NormalizeEmail(in.Email)
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, domainuser.ErrInvalidEmail
	}
	if err := uc.policy.ValidatePassword(in.Password); err != nil {
		return nil, err
	}

	if existing, err := uc.users.FindByEmail(ctx, email); err == nil && existing != nil {
		return nil, domainuser.ErrEmailTaken
	} else if err != nil && !errors.Is(err, domainuser.ErrNotFound) {
		return nil, err
	}

	hash, err := uc.hasher.Hash(in.Password)
	if err != nil {
		return nil, err
	}

	now := uc.clock.Now()
	u := domainuser.New(uc.ids.New(), email, hash, now)
	if err := uc.users.Create(ctx, u); err != nil {
		return nil, err
	}

	_ = uc.publisher.Publish(ctx, ports.Event{
		ID:         uc.ids.New(),
		Type:       EventUserRegistered,
		Version:    EventVersion,
		OccurredAt: now.UTC().Format(rfc3339),
		RoutingKey: EventUserRegistered,
		Data:       UserRegisteredData{UserID: u.ID, Email: u.Email},
	})

	return &RegisterOutput{UserID: u.ID, Email: u.Email}, nil
}

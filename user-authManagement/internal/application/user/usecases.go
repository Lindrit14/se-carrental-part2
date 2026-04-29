// Package user contains profile-management use cases.
package user

import (
	"context"

	"github.com/lindritprekaj/user-authmanagement/internal/application/auth"
	"github.com/lindritprekaj/user-authmanagement/internal/application/ports"
	"github.com/lindritprekaj/user-authmanagement/internal/domain/token"
	domainuser "github.com/lindritprekaj/user-authmanagement/internal/domain/user"
)

type Profile struct {
	ID        string
	Email     string
	Roles     []string
	Verified  bool
	CreatedAt string
}

func toProfile(u *domainuser.User) *Profile {
	roles := make([]string, 0, len(u.Roles))
	for _, r := range u.Roles {
		roles = append(roles, string(r))
	}
	return &Profile{
		ID:        u.ID,
		Email:     u.Email,
		Roles:     roles,
		Verified:  u.Verified,
		CreatedAt: u.CreatedAt.UTC().Format(rfc3339),
	}
}

const rfc3339 = "2006-01-02T15:04:05Z07:00"

// --- Get ----------------------------------------------------------------

type GetProfileUseCase struct{ users domainuser.Repository }

func NewGetProfileUseCase(users domainuser.Repository) *GetProfileUseCase {
	return &GetProfileUseCase{users: users}
}

func (uc *GetProfileUseCase) Execute(ctx context.Context, userID string) (*Profile, error) {
	u, err := uc.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return toProfile(u), nil
}

// --- Update -------------------------------------------------------------

type UpdateProfileInput struct {
	UserID string
	Email  *string
}

type UpdateProfileUseCase struct {
	users     domainuser.Repository
	publisher ports.EventPublisher
	clock     ports.Clock
	ids       ports.IDGenerator
}

func NewUpdateProfileUseCase(
	users domainuser.Repository,
	publisher ports.EventPublisher,
	clock ports.Clock,
	ids ports.IDGenerator,
) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{users: users, publisher: publisher, clock: clock, ids: ids}
}

func (uc *UpdateProfileUseCase) Execute(ctx context.Context, in UpdateProfileInput) (*Profile, error) {
	u, err := uc.users.FindByID(ctx, in.UserID)
	if err != nil {
		return nil, err
	}
	changed := make([]string, 0, 1)
	if in.Email != nil {
		newEmail := domainuser.NormalizeEmail(*in.Email)
		if newEmail != u.Email {
			u.Email = newEmail
			u.Verified = false
			changed = append(changed, "email")
		}
	}
	if len(changed) == 0 {
		return toProfile(u), nil
	}
	now := uc.clock.Now()
	u.UpdatedAt = now
	if err := uc.users.Update(ctx, u); err != nil {
		return nil, err
	}

	_ = uc.publisher.Publish(ctx, ports.Event{
		ID:         uc.ids.New(),
		Type:       auth.EventUserUpdated,
		Version:    auth.EventVersion,
		OccurredAt: now.UTC().Format(rfc3339),
		RoutingKey: auth.EventUserUpdated,
		Data:       auth.UserUpdatedData{UserID: u.ID, ChangedFields: changed},
	})
	return toProfile(u), nil
}

// --- Delete -------------------------------------------------------------

type DeleteAccountUseCase struct {
	users     domainuser.Repository
	tokens    token.Repository
	publisher ports.EventPublisher
	clock     ports.Clock
	ids       ports.IDGenerator
}

func NewDeleteAccountUseCase(
	users domainuser.Repository,
	tokens token.Repository,
	publisher ports.EventPublisher,
	clock ports.Clock,
	ids ports.IDGenerator,
) *DeleteAccountUseCase {
	return &DeleteAccountUseCase{users: users, tokens: tokens, publisher: publisher, clock: clock, ids: ids}
}

func (uc *DeleteAccountUseCase) Execute(ctx context.Context, userID string) error {
	if err := uc.tokens.RevokeAllForUser(ctx, userID); err != nil {
		return err
	}
	if err := uc.users.Delete(ctx, userID); err != nil {
		return err
	}
	now := uc.clock.Now()
	_ = uc.publisher.Publish(ctx, ports.Event{
		ID:         uc.ids.New(),
		Type:       auth.EventUserDeleted,
		Version:    auth.EventVersion,
		OccurredAt: now.UTC().Format(rfc3339),
		RoutingKey: auth.EventUserDeleted,
		Data:       auth.UserDeletedData{UserID: userID},
	})
	return nil
}

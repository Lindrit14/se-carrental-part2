// Package admin holds admin-only use cases (cross-user reads and writes).
// All public methods assume the caller has already been authorized as admin.
package admin

import (
	"context"
	"errors"

	"github.com/lindritprekaj/user-authmanagement/internal/application/auth"
	"github.com/lindritprekaj/user-authmanagement/internal/application/ports"
	"github.com/lindritprekaj/user-authmanagement/internal/domain/token"
	domainuser "github.com/lindritprekaj/user-authmanagement/internal/domain/user"
)

const rfc3339 = "2006-01-02T15:04:05Z07:00"

// ErrLastAdmin is returned when an operation would remove the very last admin.
var ErrLastAdmin = errors.New("cannot remove the last admin")

// AdminUserView is the projection used by admin endpoints.
type AdminUserView struct {
	ID        string
	Email     string
	Roles     []string
	Verified  bool
	CreatedAt string
	UpdatedAt string
}

func toView(u *domainuser.User) *AdminUserView {
	roles := make([]string, 0, len(u.Roles))
	for _, r := range u.Roles {
		roles = append(roles, string(r))
	}
	return &AdminUserView{
		ID: u.ID, Email: u.Email, Roles: roles, Verified: u.Verified,
		CreatedAt: u.CreatedAt.UTC().Format(rfc3339),
		UpdatedAt: u.UpdatedAt.UTC().Format(rfc3339),
	}
}

// --- List ---------------------------------------------------------------

type ListUsersInput struct {
	Limit  int
	Offset int
}

type ListUsersOutput struct {
	Items  []*AdminUserView
	Total  int64
	Limit  int
	Offset int
}

type ListUsersUseCase struct{ users domainuser.Repository }

func NewListUsersUseCase(users domainuser.Repository) *ListUsersUseCase {
	return &ListUsersUseCase{users: users}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context, in ListUsersInput) (*ListUsersOutput, error) {
	limit := in.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	offset := in.Offset
	if offset < 0 {
		offset = 0
	}
	items, total, err := uc.users.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	views := make([]*AdminUserView, 0, len(items))
	for _, u := range items {
		views = append(views, toView(u))
	}
	return &ListUsersOutput{Items: views, Total: total, Limit: limit, Offset: offset}, nil
}

// --- Get ----------------------------------------------------------------

type GetUserUseCase struct{ users domainuser.Repository }

func NewGetUserUseCase(users domainuser.Repository) *GetUserUseCase {
	return &GetUserUseCase{users: users}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, userID string) (*AdminUserView, error) {
	u, err := uc.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return toView(u), nil
}

// --- SetRoles -----------------------------------------------------------

type SetRolesInput struct {
	TargetUserID string
	Roles        []string
}

type SetRolesUseCase struct {
	users     domainuser.Repository
	publisher ports.EventPublisher
	clock     ports.Clock
	ids       ports.IDGenerator
}

func NewSetRolesUseCase(
	users domainuser.Repository,
	publisher ports.EventPublisher,
	clock ports.Clock,
	ids ports.IDGenerator,
) *SetRolesUseCase {
	return &SetRolesUseCase{users: users, publisher: publisher, clock: clock, ids: ids}
}

func (uc *SetRolesUseCase) Execute(ctx context.Context, in SetRolesInput) (*AdminUserView, error) {
	if len(in.Roles) == 0 {
		return nil, errors.New("at least one role is required")
	}

	target, err := uc.users.FindByID(ctx, in.TargetUserID)
	if err != nil {
		return nil, err
	}

	newRoles := dedupeRoles(in.Roles)

	// Last-admin guard: if the target currently is admin and the new role
	// list does not include admin, ensure they are not the last one.
	if target.HasRole(domainuser.RoleAdmin) && !containsRole(newRoles, domainuser.RoleAdmin) {
		count, err := uc.users.CountByRole(ctx, domainuser.RoleAdmin)
		if err != nil {
			return nil, err
		}
		if count <= 1 {
			return nil, ErrLastAdmin
		}
	}

	target.Roles = newRoles
	now := uc.clock.Now()
	target.UpdatedAt = now
	if err := uc.users.Update(ctx, target); err != nil {
		return nil, err
	}

	_ = uc.publisher.Publish(ctx, ports.Event{
		ID:         uc.ids.New(),
		Type:       auth.EventUserUpdated,
		Version:    auth.EventVersion,
		OccurredAt: now.UTC().Format(rfc3339),
		RoutingKey: auth.EventUserUpdated,
		Data:       auth.UserUpdatedData{UserID: target.ID, ChangedFields: []string{"roles"}},
	})

	return toView(target), nil
}

func dedupeRoles(in []string) []domainuser.Role {
	seen := make(map[domainuser.Role]struct{}, len(in))
	out := make([]domainuser.Role, 0, len(in))
	for _, r := range in {
		role := domainuser.Role(r)
		if _, ok := seen[role]; ok {
			continue
		}
		seen[role] = struct{}{}
		out = append(out, role)
	}
	return out
}

func containsRole(rs []domainuser.Role, target domainuser.Role) bool {
	for _, r := range rs {
		if r == target {
			return true
		}
	}
	return false
}

// --- Delete -------------------------------------------------------------

type DeleteUserUseCase struct {
	users     domainuser.Repository
	tokens    token.Repository
	publisher ports.EventPublisher
	clock     ports.Clock
	ids       ports.IDGenerator
}

func NewDeleteUserUseCase(
	users domainuser.Repository,
	tokens token.Repository,
	publisher ports.EventPublisher,
	clock ports.Clock,
	ids ports.IDGenerator,
) *DeleteUserUseCase {
	return &DeleteUserUseCase{users: users, tokens: tokens, publisher: publisher, clock: clock, ids: ids}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, targetUserID string) error {
	target, err := uc.users.FindByID(ctx, targetUserID)
	if err != nil {
		return err
	}
	if target.HasRole(domainuser.RoleAdmin) {
		count, err := uc.users.CountByRole(ctx, domainuser.RoleAdmin)
		if err != nil {
			return err
		}
		if count <= 1 {
			return ErrLastAdmin
		}
	}
	if err := uc.tokens.RevokeAllForUser(ctx, targetUserID); err != nil {
		return err
	}
	if err := uc.users.Delete(ctx, targetUserID); err != nil {
		return err
	}
	now := uc.clock.Now()
	_ = uc.publisher.Publish(ctx, ports.Event{
		ID:         uc.ids.New(),
		Type:       auth.EventUserDeleted,
		Version:    auth.EventVersion,
		OccurredAt: now.UTC().Format(rfc3339),
		RoutingKey: auth.EventUserDeleted,
		Data:       auth.UserDeletedData{UserID: targetUserID},
	})
	return nil
}

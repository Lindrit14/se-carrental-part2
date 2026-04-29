package admin

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lindritprekaj/user-authmanagement/internal/application/ports"
	domainuser "github.com/lindritprekaj/user-authmanagement/internal/domain/user"
)

// --- Test doubles -------------------------------------------------------

type fakeRepo struct {
	users     map[string]*domainuser.User
	adminCnt  int64
	updateErr error
}

func newFakeRepo() *fakeRepo { return &fakeRepo{users: map[string]*domainuser.User{}} }

func (f *fakeRepo) Create(_ context.Context, u *domainuser.User) error {
	f.users[u.ID] = u
	return nil
}

func (f *fakeRepo) FindByID(_ context.Context, id string) (*domainuser.User, error) {
	u, ok := f.users[id]
	if !ok {
		return nil, domainuser.ErrNotFound
	}
	return u, nil
}

func (f *fakeRepo) FindByEmail(context.Context, string) (*domainuser.User, error) {
	return nil, domainuser.ErrNotFound
}

func (f *fakeRepo) Update(_ context.Context, u *domainuser.User) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	f.users[u.ID] = u
	return nil
}

func (f *fakeRepo) Delete(context.Context, string) error  { return nil }
func (f *fakeRepo) List(context.Context, int, int) ([]*domainuser.User, int64, error) {
	return nil, 0, nil
}
func (f *fakeRepo) CountByRole(context.Context, domainuser.Role) (int64, error) {
	return f.adminCnt, nil
}

type fakePublisher struct{ events []ports.Event }

func (f *fakePublisher) Publish(_ context.Context, e ports.Event) error {
	f.events = append(f.events, e)
	return nil
}

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }

type fixedIDs struct{ id string }

func (f fixedIDs) New() string { return f.id }

// --- Tests --------------------------------------------------------------

func TestSetRoles_GrantsAdminAndPublishesEvent(t *testing.T) {
	repo := newFakeRepo()
	now := time.Now().UTC()
	repo.users["u1"] = domainuser.New("u1", "user@example.com", "hash", now)
	repo.adminCnt = 1

	pub := &fakePublisher{}
	uc := NewSetRolesUseCase(repo, pub, fixedClock{t: now}, fixedIDs{id: "evt"})

	out, err := uc.Execute(context.Background(), SetRolesInput{
		TargetUserID: "u1",
		Roles:        []string{"user", "admin"},
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if !contains(out.Roles, "admin") {
		t.Fatalf("expected admin role in output, got %v", out.Roles)
	}
	if len(pub.events) != 1 {
		t.Fatalf("expected exactly 1 event, got %d", len(pub.events))
	}
	if pub.events[0].Type != "user.updated" {
		t.Errorf("expected event type user.updated, got %s", pub.events[0].Type)
	}
}

func TestSetRoles_RejectsRemovingLastAdmin(t *testing.T) {
	repo := newFakeRepo()
	now := time.Now().UTC()
	admin := domainuser.New("u1", "admin@example.com", "hash", now)
	admin.Roles = []domainuser.Role{domainuser.RoleUser, domainuser.RoleAdmin}
	repo.users["u1"] = admin
	repo.adminCnt = 1

	uc := NewSetRolesUseCase(repo, &fakePublisher{}, fixedClock{t: now}, fixedIDs{id: "evt"})

	_, err := uc.Execute(context.Background(), SetRolesInput{
		TargetUserID: "u1",
		Roles:        []string{"user"},
	})
	if !errors.Is(err, ErrLastAdmin) {
		t.Fatalf("expected ErrLastAdmin, got %v", err)
	}
}

func TestSetRoles_AllowsRemovingAdminWhenNotLast(t *testing.T) {
	repo := newFakeRepo()
	now := time.Now().UTC()
	admin := domainuser.New("u1", "admin@example.com", "hash", now)
	admin.Roles = []domainuser.Role{domainuser.RoleUser, domainuser.RoleAdmin}
	repo.users["u1"] = admin
	repo.adminCnt = 2 // there is at least one other admin

	uc := NewSetRolesUseCase(repo, &fakePublisher{}, fixedClock{t: now}, fixedIDs{id: "evt"})

	out, err := uc.Execute(context.Background(), SetRolesInput{
		TargetUserID: "u1",
		Roles:        []string{"user"},
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if contains(out.Roles, "admin") {
		t.Fatalf("expected admin role to be removed, got %v", out.Roles)
	}
}

func contains(rs []string, target string) bool {
	for _, r := range rs {
		if r == target {
			return true
		}
	}
	return false
}

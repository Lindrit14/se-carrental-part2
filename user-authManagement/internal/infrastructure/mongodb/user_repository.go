package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	domainuser "github.com/lindritprekaj/user-authmanagement/internal/domain/user"
)

type userDoc struct {
	ID           string    `bson:"_id"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password_hash"`
	Roles        []string  `bson:"roles"`
	Verified     bool      `bson:"verified"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

func toUserDoc(u *domainuser.User) userDoc {
	roles := make([]string, 0, len(u.Roles))
	for _, r := range u.Roles {
		roles = append(roles, string(r))
	}
	return userDoc{
		ID: u.ID, Email: u.Email, PasswordHash: u.PasswordHash,
		Roles: roles, Verified: u.Verified,
		CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
	}
}

func fromUserDoc(d userDoc) *domainuser.User {
	roles := make([]domainuser.Role, 0, len(d.Roles))
	for _, r := range d.Roles {
		roles = append(roles, domainuser.Role(r))
	}
	return &domainuser.User{
		ID: d.ID, Email: d.Email, PasswordHash: d.PasswordHash,
		Roles: roles, Verified: d.Verified,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

type UserRepository struct{ coll *mongo.Collection }

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{coll: db.Collection(CollectionUsers)}
}

func (r *UserRepository) Create(ctx context.Context, u *domainuser.User) error {
	_, err := r.coll.InsertOne(ctx, toUserDoc(u))
	if err != nil && mongo.IsDuplicateKeyError(err) {
		return domainuser.ErrEmailTaken
	}
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domainuser.User, error) {
	var d userDoc
	if err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&d); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domainuser.ErrNotFound
		}
		return nil, err
	}
	return fromUserDoc(d), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	var d userDoc
	if err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&d); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domainuser.ErrNotFound
		}
		return nil, err
	}
	return fromUserDoc(d), nil
}

func (r *UserRepository) Update(ctx context.Context, u *domainuser.User) error {
	res, err := r.coll.ReplaceOne(ctx, bson.M{"_id": u.ID}, toUserDoc(u))
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domainuser.ErrEmailTaken
		}
		return err
	}
	if res.MatchedCount == 0 {
		return domainuser.ErrNotFound
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return domainuser.ErrNotFound
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*domainuser.User, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	total, err := r.coll.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}
	cur, err := r.coll.Find(ctx, bson.D{}, options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset)),
	)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var docs []userDoc
	if err := cur.All(ctx, &docs); err != nil {
		return nil, 0, err
	}
	out := make([]*domainuser.User, 0, len(docs))
	for _, d := range docs {
		out = append(out, fromUserDoc(d))
	}
	return out, total, nil
}

func (r *UserRepository) CountByRole(ctx context.Context, role domainuser.Role) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{"roles": string(role)})
}

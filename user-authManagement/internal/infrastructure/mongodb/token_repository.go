package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/lindritprekaj/user-authmanagement/internal/domain/token"
)

type refreshDoc struct {
	ID         string     `bson:"_id"`
	UserID     string     `bson:"user_id"`
	TokenHash  string     `bson:"token_hash"`
	IssuedAt   time.Time  `bson:"issued_at"`
	ExpiresAt  time.Time  `bson:"expires_at"`
	RevokedAt  *time.Time `bson:"revoked_at,omitempty"`
	ReplacedBy string     `bson:"replaced_by,omitempty"`
}

func toRefreshDoc(t *token.RefreshToken) refreshDoc {
	return refreshDoc{
		ID: t.ID, UserID: t.UserID, TokenHash: t.TokenHash,
		IssuedAt: t.IssuedAt, ExpiresAt: t.ExpiresAt,
		RevokedAt: t.RevokedAt, ReplacedBy: t.ReplacedBy,
	}
}

func fromRefreshDoc(d refreshDoc) *token.RefreshToken {
	return &token.RefreshToken{
		ID: d.ID, UserID: d.UserID, TokenHash: d.TokenHash,
		IssuedAt: d.IssuedAt, ExpiresAt: d.ExpiresAt,
		RevokedAt: d.RevokedAt, ReplacedBy: d.ReplacedBy,
	}
}

type RefreshTokenRepository struct{ coll *mongo.Collection }

func NewRefreshTokenRepository(db *mongo.Database) *RefreshTokenRepository {
	return &RefreshTokenRepository{coll: db.Collection(CollectionRefreshTokens)}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, t *token.RefreshToken) error {
	_, err := r.coll.InsertOne(ctx, toRefreshDoc(t))
	return err
}

func (r *RefreshTokenRepository) FindByHash(ctx context.Context, tokenHash string) (*token.RefreshToken, error) {
	var d refreshDoc
	if err := r.coll.FindOne(ctx, bson.M{"token_hash": tokenHash}).Decode(&d); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, token.ErrNotFound
		}
		return nil, err
	}
	return fromRefreshDoc(d), nil
}

func (r *RefreshTokenRepository) Update(ctx context.Context, t *token.RefreshToken) error {
	res, err := r.coll.ReplaceOne(ctx, bson.M{"_id": t.ID}, toRefreshDoc(t))
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return token.ErrNotFound
	}
	return nil
}

func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	now := time.Now().UTC()
	_, err := r.coll.UpdateMany(ctx,
		bson.M{"user_id": userID, "revoked_at": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{"revoked_at": now}},
	)
	return err
}

// --- Password reset --------------------------------------------------

type passwordResetDoc struct {
	ID        string     `bson:"_id"`
	UserID    string     `bson:"user_id"`
	TokenHash string     `bson:"token_hash"`
	ExpiresAt time.Time  `bson:"expires_at"`
	UsedAt    *time.Time `bson:"used_at,omitempty"`
}

type PasswordResetRepository struct{ coll *mongo.Collection }

func NewPasswordResetRepository(db *mongo.Database) *PasswordResetRepository {
	return &PasswordResetRepository{coll: db.Collection(CollectionPasswordResets)}
}

func (r *PasswordResetRepository) Create(ctx context.Context, p *token.PasswordReset) error {
	_, err := r.coll.InsertOne(ctx, passwordResetDoc{
		ID: p.ID, UserID: p.UserID, TokenHash: p.TokenHash,
		ExpiresAt: p.ExpiresAt, UsedAt: p.UsedAt,
	})
	return err
}

func (r *PasswordResetRepository) FindByHash(ctx context.Context, tokenHash string) (*token.PasswordReset, error) {
	var d passwordResetDoc
	if err := r.coll.FindOne(ctx, bson.M{"token_hash": tokenHash}).Decode(&d); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, token.ErrNotFound
		}
		return nil, err
	}
	return &token.PasswordReset{
		ID: d.ID, UserID: d.UserID, TokenHash: d.TokenHash,
		ExpiresAt: d.ExpiresAt, UsedAt: d.UsedAt,
	}, nil
}

func (r *PasswordResetRepository) MarkUsed(ctx context.Context, id string, usedAt time.Time) error {
	res, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"used_at": usedAt}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return token.ErrNotFound
	}
	return nil
}

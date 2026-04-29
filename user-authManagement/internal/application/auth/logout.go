package auth

import (
	"context"

	"github.com/lindritprekaj/user-authmanagement/internal/application/ports"
	"github.com/lindritprekaj/user-authmanagement/internal/domain/token"
)

type LogoutInput struct {
	RefreshToken string
}

type LogoutUseCase struct {
	tokens   token.Repository
	tokenSvc ports.TokenService
	clock    ports.Clock
}

func NewLogoutUseCase(tokens token.Repository, tokenSvc ports.TokenService, clock ports.Clock) *LogoutUseCase {
	return &LogoutUseCase{tokens: tokens, tokenSvc: tokenSvc, clock: clock}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, in LogoutInput) error {
	hash := uc.tokenSvc.HashRefreshToken(in.RefreshToken)
	rt, err := uc.tokens.FindByHash(ctx, hash)
	if err != nil {
		return err
	}
	rt.Revoke(uc.clock.Now(), "")
	return uc.tokens.Update(ctx, rt)
}

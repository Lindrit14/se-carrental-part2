package handlers

import (
	"log/slog"
	"net/http"

	usercases "github.com/lindritprekaj/user-authmanagement/internal/application/user"
	"github.com/lindritprekaj/user-authmanagement/internal/interfaces/http/dto"
	"github.com/lindritprekaj/user-authmanagement/internal/interfaces/http/middleware"
)

type UserHandler struct {
	logger      *slog.Logger
	getProfile  *usercases.GetProfileUseCase
	updateProf  *usercases.UpdateProfileUseCase
	deleteAcct  *usercases.DeleteAccountUseCase
}

func NewUserHandler(
	logger *slog.Logger,
	getProfile *usercases.GetProfileUseCase,
	updateProf *usercases.UpdateProfileUseCase,
	deleteAcct *usercases.DeleteAccountUseCase,
) *UserHandler {
	return &UserHandler{
		logger: logger, getProfile: getProfile,
		updateProf: updateProf, deleteAcct: deleteAcct,
	}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UserIDFromContext(r.Context())
	p, err := h.getProfile.Execute(r.Context(), uid)
	if err != nil {
		writeDomainError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusOK, toProfileResponse(p))
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UserIDFromContext(r.Context())
	var req dto.UpdateProfileRequest
	if err := decodeAndValidate(w, r, &req); err != nil {
		return
	}
	p, err := h.updateProf.Execute(r.Context(), usercases.UpdateProfileInput{
		UserID: uid, Email: req.Email,
	})
	if err != nil {
		writeDomainError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusOK, toProfileResponse(p))
}

func (h *UserHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UserIDFromContext(r.Context())
	if err := h.deleteAcct.Execute(r.Context(), uid); err != nil {
		writeDomainError(w, r, h.logger, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toProfileResponse(p *usercases.Profile) dto.ProfileResponse {
	return dto.ProfileResponse{
		ID: p.ID, Email: p.Email, Roles: p.Roles,
		Verified: p.Verified, CreatedAt: p.CreatedAt,
	}
}

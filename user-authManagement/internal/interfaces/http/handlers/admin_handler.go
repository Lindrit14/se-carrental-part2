package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/lindritprekaj/user-authmanagement/internal/application/admin"
	"github.com/lindritprekaj/user-authmanagement/internal/interfaces/http/dto"
)

type AdminHandler struct {
	logger     *slog.Logger
	listUsers  *admin.ListUsersUseCase
	getUser    *admin.GetUserUseCase
	setRoles   *admin.SetRolesUseCase
	deleteUser *admin.DeleteUserUseCase
}

func NewAdminHandler(
	logger *slog.Logger,
	listUsers *admin.ListUsersUseCase,
	getUser *admin.GetUserUseCase,
	setRoles *admin.SetRolesUseCase,
	deleteUser *admin.DeleteUserUseCase,
) *AdminHandler {
	return &AdminHandler{
		logger: logger, listUsers: listUsers, getUser: getUser,
		setRoles: setRoles, deleteUser: deleteUser,
	}
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	out, err := h.listUsers.Execute(r.Context(), admin.ListUsersInput{Limit: limit, Offset: offset})
	if err != nil {
		writeDomainError(w, r, h.logger, err)
		return
	}
	items := make([]dto.AdminUserResponse, 0, len(out.Items))
	for _, v := range out.Items {
		items = append(items, toAdminUserResponse(v))
	}
	writeJSON(w, http.StatusOK, dto.AdminUserListResponse{
		Items: items, Total: out.Total, Limit: out.Limit, Offset: out.Offset,
	})
}

func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	v, err := h.getUser.Execute(r.Context(), id)
	if err != nil {
		writeDomainError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusOK, toAdminUserResponse(v))
}

func (h *AdminHandler) SetRoles(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.AdminSetRolesRequest
	if err := decodeAndValidate(w, r, &req); err != nil {
		return
	}
	v, err := h.setRoles.Execute(r.Context(), admin.SetRolesInput{
		TargetUserID: id, Roles: req.Roles,
	})
	if err != nil {
		writeDomainError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusOK, toAdminUserResponse(v))
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.deleteUser.Execute(r.Context(), id); err != nil {
		writeDomainError(w, r, h.logger, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toAdminUserResponse(v *admin.AdminUserView) dto.AdminUserResponse {
	return dto.AdminUserResponse{
		ID: v.ID, Email: v.Email, Roles: v.Roles, Verified: v.Verified,
		CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt,
	}
}

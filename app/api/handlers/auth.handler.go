package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"orijinplus/app/models"
	"orijinplus/app/services"
	"orijinplus/utils/authtoken"
	"orijinplus/utils/faulterr"
	"time"
)

type AuthHandler struct {
	services *services.Services
}

func NewAuthHandler(s *services.Services) *AuthHandler {
	return &AuthHandler{s}
}

// GenerateToken generates a jwt token and sets cookie
func (h *AuthHandler) GenerateToken(w http.ResponseWriter, ctx context.Context, auther *models.Auther) (*http.Cookie, *faulterr.FaultErr) {
	tokenPayload, err := authtoken.Generate(auther)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    tokenPayload.TokenString,
		Expires:  tokenPayload.ExpiresAt,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	return cookie, nil
}

// ListPermissions Handler
func (h *AuthHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	result := models.ListPermissions()

	response := ResponseBody{
		Data:       result,
		Message:    "Permissions list",
		StatusCode: http.StatusOK,
	}

	RestResponse(w, r, response.StatusCode, response)
}

// LoginAdmin return jwt token of superadmin and prevents other users
func (h *AuthHandler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := &models.LoginRequest{}
	if decodeErr := json.NewDecoder(r.Body).Decode(request); decodeErr != nil {
		err := faulterr.NewUnprocessableEntityError("Invalid JSON request")
		RestResponse(w, r, err.Status, err)
		return
	}
	defer r.Body.Close()

	if err := h.services.AuthService.ValidateLoginRequest(*request); err != nil {
		RestResponse(w, r, err.Status, err)
		return
	}

	auther, err := h.services.AuthService.Login(r.Context(), request)
	if err != nil {
		RestResponse(w, r, err.Status, err)
		return
	}
	if !auther.IsAdmin {
		err := faulterr.NewNotFoundError("user not found")
		RestResponse(w, r, err.Status, err)
		return
	}

	cookieToken, err := h.GenerateToken(w, ctx, auther)
	if err != nil {
		RestResponse(w, r, err.Status, err)
		return
	}

	authData := AuthData{
		CookieToken: cookieToken,
		Auther:      auther,
		Token:       cookieToken.Value,
	}

	response := ResponseBody{
		Data:       authData,
		Message:    "Admin logged in successfully",
		StatusCode: http.StatusAccepted,
	}

	RestResponse(w, r, response.StatusCode, response)
}

// LoginMember return jwt token of member and prevents other users
func (h *AuthHandler) LoginMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := &models.LoginRequest{}
	if decodeErr := json.NewDecoder(r.Body).Decode(request); decodeErr != nil {
		err := faulterr.NewUnprocessableEntityError("Invalid JSON request")
		RestResponse(w, r, err.Status, err)
		return
	}
	defer r.Body.Close()

	if err := h.services.AuthService.ValidateLoginRequest(*request); err != nil {
		RestResponse(w, r, err.Status, err)
		return
	}

	auther, err := h.services.AuthService.Login(r.Context(), request)
	if err != nil {
		RestResponse(w, r, err.Status, err)
		return
	}
	if !auther.IsMember {
		err := faulterr.NewNotFoundError("user not found")
		RestResponse(w, r, err.Status, err)
		return
	}

	cookieToken, err := h.GenerateToken(w, ctx, auther)
	if err != nil {
		RestResponse(w, r, err.Status, err)
		return
	}

	authData := AuthData{
		CookieToken: cookieToken,
		Auther:      auther,
		Token:       cookieToken.Value,
	}

	response := ResponseBody{
		Data:       authData,
		Message:    "Member logged in successfully",
		StatusCode: http.StatusAccepted,
	}

	RestResponse(w, r, response.StatusCode, response)
}

// LoginCustomer return jwt token of customer and prevents other users
func (h *AuthHandler) LoginCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := &models.LoginRequest{}
	if decodeErr := json.NewDecoder(r.Body).Decode(request); decodeErr != nil {
		err := faulterr.NewUnprocessableEntityError("Invalid JSON request")
		RestResponse(w, r, err.Status, err)
		return
	}
	defer r.Body.Close()

	if err := h.services.AuthService.ValidateLoginRequest(*request); err != nil {
		RestResponse(w, r, err.Status, err)
		return
	}

	auther, err := h.services.AuthService.Login(r.Context(), request)
	if err != nil {
		RestResponse(w, r, err.Status, err)
		return
	}
	if !auther.IsCustomer {
		err := faulterr.NewNotFoundError("user not found")
		RestResponse(w, r, err.Status, err)
		return
	}

	cookieToken, err := h.GenerateToken(w, ctx, auther)
	if err != nil {
		RestResponse(w, r, err.Status, err)
		return
	}

	authData := AuthData{
		CookieToken: cookieToken,
		Auther:      auther,
		Token:       cookieToken.Value,
	}

	response := ResponseBody{
		Data:       authData,
		Message:    "Customer logged in successfully",
		StatusCode: http.StatusAccepted,
	}

	RestResponse(w, r, response.StatusCode, response)
}

// RegisterAdmin Handler
func (h *AuthHandler) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	request := &models.SuperAdminRequest{}
	if decodeErr := json.NewDecoder(r.Body).Decode(&request); decodeErr != nil {
		err := faulterr.NewUnprocessableEntityError("Invalid JSON request")
		RestResponse(w, r, err.Status, err)
		return
	}
	defer r.Body.Close()

	result, err := h.services.AuthService.RegisterAdmin(r.Context(), *request)
	if err != nil {
		RestResponse(w, r, err.Status, err)
		return
	}

	response := ResponseBody{
		Data:       result,
		Message:    "Admin registered successfully",
		StatusCode: http.StatusCreated,
	}

	RestResponse(w, r, response.StatusCode, response)
}

// RegisterMember Handler
func (h *AuthHandler) RegisterMember(w http.ResponseWriter, r *http.Request) {
	request := models.MemberRequest{}
	if decodeErr := json.NewDecoder(r.Body).Decode(&request); decodeErr != nil {
		err := faulterr.NewUnprocessableEntityError("Invalid JSON request")
		RestResponse(w, r, err.Status, err)
		return
	}
	defer r.Body.Close()

	result, err := h.services.AuthService.RegisterMember(r.Context(), request)
	if err != nil {
		RestResponse(w, r, err.Status, err)
		return
	}

	response := ResponseBody{
		Data:       result,
		Message:    "Member registered successfully",
		StatusCode: http.StatusCreated,
	}

	RestResponse(w, r, response.StatusCode, response)
}

// RegisterCustomer Handler
func (h *AuthHandler) RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	request := models.CustomerRequest{}
	if decodeErr := json.NewDecoder(r.Body).Decode(&request); decodeErr != nil {
		err := faulterr.NewUnprocessableEntityError("Invalid JSON request")
		RestResponse(w, r, err.Status, err)
		return
	}
	defer r.Body.Close()

	result, err := h.services.AuthService.RegisterCustomer(r.Context(), request)
	if err != nil {
		RestResponse(w, r, err.Status, err)
		return
	}

	response := ResponseBody{
		Data:       result,
		Message:    "Customer registered successfully",
		StatusCode: http.StatusCreated,
	}

	RestResponse(w, r, response.StatusCode, response)
}

// RegisterOrganization Handler
func (h *AuthHandler) RegisterOrganization(w http.ResponseWriter, r *http.Request) {
	request := models.OrganizationRequest{}
	if decodeErr := json.NewDecoder(r.Body).Decode(&request); decodeErr != nil {
		err := faulterr.NewUnprocessableEntityError("Invalid JSON request")
		RestResponse(w, r, err.Status, err)
		return
	}
	defer r.Body.Close()

	result, err := h.services.AuthService.RegisterOrganization(r.Context(), request)
	if err != nil {
		RestResponse(w, r, err.Status, err)
		return
	}

	response := ResponseBody{
		Data:       result,
		Message:    "Organization registered successfully",
		StatusCode: http.StatusCreated,
	}

	RestResponse(w, r, response.StatusCode, response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w,
		&http.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		})

	response := ResponseBody{
		Data:       nil,
		Message:    "user logout",
		StatusCode: http.StatusCreated,
	}

	RestResponse(w, r, response.StatusCode, response)
}

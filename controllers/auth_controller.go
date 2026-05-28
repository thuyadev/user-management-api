package controllers

import (
	"errors"
	"net/http"

	"user-management-api/auth"
	"user-management-api/middleware"
	"user-management-api/models"
	"user-management-api/services"
	"user-management-api/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Login godoc
// @Summary      Login
// @Description  Authenticate with email and password. Returns JWT token, role, and permissions.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.LoginRequest  true  "Login credentials"
// @Success      200   {object}  models.SwaggerAuthTokenSuccessResponse
// @Failure      401   {object}  models.SwaggerErrorResponse
// @Failure      422   {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	resp, err := ctrl.authService.Login(req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			utils.Unauthorized(c, "Invalid email or password")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Login successful", resp)
}

// Register godoc
// @Summary      Register
// @Description  Create a new user account (role: user). Returns JWT token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.RegisterRequest  true  "Registration data"
// @Success      201   {object}  models.SwaggerAuthTokenSuccessResponse
// @Failure      409   {object}  models.SwaggerErrorResponse
// @Failure      422   {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	resp, err := ctrl.authService.Register(req)
	if err != nil {
		if errors.Is(err, services.ErrEmailTaken) {
			utils.Error(c, http.StatusConflict, "Email already registered", nil)
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, "Registration successful", resp)
}

// Me godoc
// @Summary      Current user profile
// @Description  Returns authenticated user id, email, role, and permissions.
// @Tags         Auth
// @Produce      json
// @Success      200  {object}  models.SwaggerMeResponse
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/auth/me [get]
func (ctrl *AuthController) Me(c *gin.Context) {
	role := c.GetString(middleware.ContextUserRoleKey)
	utils.Success(c, http.StatusOK, "Profile retrieved", gin.H{
		"user_id":     middleware.GetUserID(c),
		"email":       c.GetString(middleware.ContextUserEmailKey),
		"role":        role,
		"permissions": auth.PermissionsForRole(role),
	})
}

// Roles godoc
// @Summary      List roles and permissions
// @Description  Returns all supported roles and their permission lists.
// @Tags         Auth
// @Produce      json
// @Success      200  {object}  models.SwaggerRolesResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/auth/roles [get]
func (ctrl *AuthController) Roles(c *gin.Context) {
	utils.Success(c, http.StatusOK, "Roles and permissions", auth.AllRoles())
}

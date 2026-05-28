package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"user-management-api/middleware"
	"user-management-api/models"
	"user-management-api/services"
	"user-management-api/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

// List godoc
// @Summary      List users
// @Description  Paginated user list. Requires users.manage permission (admin only).
// @Tags         Users
// @Produce      json
// @Param        page      query  int     false  "Page number"       default(1)
// @Param        per_page  query  int     false  "Items per page"    default(10)
// @Param        search    query  string  false  "Search by name or email"
// @Success      200  {object}  models.SwaggerPaginatedUsersResponse
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      403  {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/users [get]
func (ctrl *UserController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	users, total, err := ctrl.userService.List(page, perPage, search)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, users, total, page, perPage)
}

// Get godoc
// @Summary      Get user
// @Description  Get a single user by ID. Requires users.manage permission.
// @Tags         Users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.SwaggerUserSuccessResponse
// @Failure      404  {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/users/{id} [get]
func (ctrl *UserController) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	user, err := ctrl.userService.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			utils.NotFound(c, "User not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "User retrieved", user)
}

// Create godoc
// @Summary      Create user
// @Description  Create a new user. Requires users.manage permission.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body  body      models.CreateUserRequest  true  "User data"
// @Success      201   {object}  models.SwaggerUserSuccessResponse
// @Failure      409   {object}  models.SwaggerErrorResponse
// @Failure      422   {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/users [post]
func (ctrl *UserController) Create(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	user, err := ctrl.userService.Create(req, middleware.GetUserID(c))
	if err != nil {
		if errors.Is(err, services.ErrEmailTaken) {
			utils.Error(c, http.StatusConflict, "Email already registered", nil)
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, "User created", user)
}

// Update godoc
// @Summary      Update user
// @Description  Update an existing user. Requires users.manage permission.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      int                       true  "User ID"
// @Param        body  body      models.UpdateUserRequest  true  "User data"
// @Success      200   {object}  models.SwaggerUserSuccessResponse
// @Failure      404   {object}  models.SwaggerErrorResponse
// @Failure      409   {object}  models.SwaggerErrorResponse
// @Failure      422   {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/users/{id} [put]
func (ctrl *UserController) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	user, err := ctrl.userService.Update(id, req, middleware.GetUserID(c))
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			utils.NotFound(c, "User not found")
			return
		}
		if errors.Is(err, services.ErrEmailTaken) {
			utils.Error(c, http.StatusConflict, "Email already registered", nil)
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "User updated", user)
}

// Delete godoc
// @Summary      Delete user
// @Description  Soft-delete a user. Requires users.manage permission.
// @Tags         Users
// @Produce      json
// @Param        id   path  int  true  "User ID"
// @Success      200  {object}  models.SwaggerSuccessResponse
// @Failure      404  {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/users/{id} [delete]
func (ctrl *UserController) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	if err := ctrl.userService.Delete(id, middleware.GetUserID(c)); err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			utils.NotFound(c, "User not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "User deleted", nil)
}

func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(id), err
}

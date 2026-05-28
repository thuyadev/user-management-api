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

type CategoryController struct {
	categoryService services.CategoryService
	aiService       services.AIService
}

func NewCategoryController(categoryService services.CategoryService, aiService services.AIService) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
		aiService:       aiService,
	}
}

// List godoc
// @Summary      List categories
// @Description  Paginated category list. Requires categories.view permission.
// @Tags         Categories
// @Produce      json
// @Param        page      query  int     false  "Page number"     default(1)
// @Param        per_page  query  int     false  "Items per page"  default(10)
// @Param        search    query  string  false  "Search by name"
// @Success      200  {object}  models.SwaggerPaginatedCategoriesResponse
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      403  {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/categories [get]
func (ctrl *CategoryController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	categories, total, err := ctrl.categoryService.List(page, perPage, search)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, categories, total, page, perPage)
}

// Get godoc
// @Summary      Get category
// @Description  Get a single category by ID. Requires categories.view permission.
// @Tags         Categories
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  models.SwaggerCategorySuccessResponse
// @Failure      404  {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/categories/{id} [get]
func (ctrl *CategoryController) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid category ID", nil)
		return
	}

	category, err := ctrl.categoryService.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			utils.NotFound(c, "Category not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Category retrieved", category)
}

// Create godoc
// @Summary      Create category
// @Description  Create a new category. Requires categories.manage permission (admin only).
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        body  body      models.CreateCategoryRequest  true  "Category data"
// @Success      201   {object}  models.SwaggerCategorySuccessResponse
// @Failure      422   {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/categories [post]
func (ctrl *CategoryController) Create(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	category, err := ctrl.categoryService.Create(req, middleware.GetUserID(c))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, "Category created", category)
}

// Update godoc
// @Summary      Update category
// @Description  Update an existing category. Requires categories.manage permission.
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id    path      int                           true  "Category ID"
// @Param        body  body      models.UpdateCategoryRequest  true  "Category data"
// @Success      200   {object}  models.SwaggerCategorySuccessResponse
// @Failure      404   {object}  models.SwaggerErrorResponse
// @Failure      422   {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/categories/{id} [put]
func (ctrl *CategoryController) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid category ID", nil)
		return
	}

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	category, err := ctrl.categoryService.Update(id, req, middleware.GetUserID(c))
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			utils.NotFound(c, "Category not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Category updated", category)
}

// Delete godoc
// @Summary      Delete category
// @Description  Soft-delete a category. Requires categories.manage permission.
// @Tags         Categories
// @Produce      json
// @Param        id   path  int  true  "Category ID"
// @Success      200  {object}  models.SwaggerSuccessResponse
// @Failure      404  {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/categories/{id} [delete]
func (ctrl *CategoryController) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid category ID", nil)
		return
	}

	if err := ctrl.categoryService.Delete(id, middleware.GetUserID(c)); err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			utils.NotFound(c, "Category not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Category deleted", nil)
}

type AISuggestCategoryRequest struct {
	Keywords string `json:"keywords" binding:"required,min=2"`
}

// SuggestName godoc
// @Summary      AI suggest category name
// @Description  Suggest a category name from keywords. Requires categories.manage permission.
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        body  body      AISuggestCategoryRequest  true  "Keywords"
// @Success      200   {object}  models.SwaggerAISuggestResponse
// @Failure      422   {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/categories/ai/suggest [post]
func (ctrl *CategoryController) SuggestName(c *gin.Context) {
	var req AISuggestCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	name, err := ctrl.aiService.SuggestCategoryName(req.Keywords)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Category name suggested", gin.H{"name": name})
}

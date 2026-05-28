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

type ProductController struct {
	productService services.ProductService
	aiService      services.AIService
}

func NewProductController(productService services.ProductService, aiService services.AIService) *ProductController {
	return &ProductController{
		productService: productService,
		aiService:      aiService,
	}
}

// List godoc
// @Summary      List products
// @Description  Paginated product list. Requires products.view permission.
// @Tags         Products
// @Produce      json
// @Param        page         query  int     false  "Page number"       default(1)
// @Param        per_page     query  int     false  "Items per page"    default(10)
// @Param        search       query  string  false  "Search by name or description"
// @Param        category_id  query  int     false  "Filter by category ID"
// @Success      200  {object}  models.SwaggerPaginatedProductsResponse
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      403  {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/products [get]
func (ctrl *ProductController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("search")
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 64)

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	products, total, err := ctrl.productService.List(page, perPage, search, uint(categoryID))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, products, total, page, perPage)
}

// Get godoc
// @Summary      Get product
// @Description  Get a single product by ID. Requires products.view permission.
// @Tags         Products
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  models.SwaggerProductSuccessResponse
// @Failure      404  {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/products/{id} [get]
func (ctrl *ProductController) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	product, err := ctrl.productService.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			utils.NotFound(c, "Product not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Product retrieved", product)
}

// Create godoc
// @Summary      Create product
// @Description  Create a new product. Requires products.manage permission (admin only).
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        body  body      models.CreateProductRequest  true  "Product data"
// @Success      201   {object}  models.SwaggerProductSuccessResponse
// @Failure      404   {object}  models.SwaggerErrorResponse
// @Failure      422   {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/products [post]
func (ctrl *ProductController) Create(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	product, err := ctrl.productService.Create(req, middleware.GetUserID(c))
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			utils.NotFound(c, "Category not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, "Product created", product)
}

// Update godoc
// @Summary      Update product
// @Description  Update an existing product. Requires products.manage permission.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id    path      int                          true  "Product ID"
// @Param        body  body      models.UpdateProductRequest  true  "Product data"
// @Success      200   {object}  models.SwaggerProductSuccessResponse
// @Failure      404   {object}  models.SwaggerErrorResponse
// @Failure      422   {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/products/{id} [put]
func (ctrl *ProductController) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	product, err := ctrl.productService.Update(id, req, middleware.GetUserID(c))
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			utils.NotFound(c, "Product not found")
			return
		}
		if errors.Is(err, services.ErrCategoryNotFound) {
			utils.NotFound(c, "Category not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Product updated", product)
}

// Delete godoc
// @Summary      Delete product
// @Description  Soft-delete a product. Requires products.manage permission.
// @Tags         Products
// @Produce      json
// @Param        id   path  int  true  "Product ID"
// @Success      200  {object}  models.SwaggerSuccessResponse
// @Failure      404  {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/products/{id} [delete]
func (ctrl *ProductController) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	if err := ctrl.productService.Delete(id, middleware.GetUserID(c)); err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			utils.NotFound(c, "Product not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Product deleted", nil)
}

type AIGenerateDescriptionRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Category string `json:"category" binding:"required,min=2"`
}

// GenerateDescription godoc
// @Summary      AI generate product description
// @Description  Generate a product description from name and category. Requires products.manage permission.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        body  body      AIGenerateDescriptionRequest  true  "Product info"
// @Success      200   {object}  models.SwaggerAIDescriptionResponse
// @Failure      422   {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/products/ai/description [post]
func (ctrl *ProductController) GenerateDescription(c *gin.Context) {
	var req AIGenerateDescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	description, err := ctrl.aiService.GenerateProductDescription(req.Name, req.Category)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Description generated", gin.H{"description": description})
}

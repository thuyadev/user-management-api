package controllers

import (
	"net/http"
	"strconv"

	"user-management-api/services"
	"user-management-api/utils"

	"github.com/gin-gonic/gin"
)

type LogController struct {
	logService services.LogService
}

func NewLogController(logService services.LogService) *LogController {
	return &LogController{logService: logService}
}

// List godoc
// @Summary      List activity logs
// @Description  Paginated user activity logs from MongoDB. Requires logs.view permission (admin only).
// @Tags         Logs
// @Produce      json
// @Param        page      query  int  false  "Page number"     default(1)
// @Param        per_page  query  int  false  "Items per page"  default(10)
// @Param        user_id   query  int  false  "Filter by user ID"
// @Success      200  {object}  models.SwaggerPaginatedLogsResponse
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      403  {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/logs [get]
func (ctrl *LogController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	logs, total, err := ctrl.logService.List(c.Request.Context(), uint(userID), page, perPage)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, logs, total, page, perPage)
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Public health check endpoint. No API key required.
// @Tags         Health
// @Produce      json
// @Success      200  {object}  models.SwaggerHealthResponse
// @Router       /health [get]
func HealthCheck(c *gin.Context) {
	utils.Success(c, http.StatusOK, "API is running", gin.H{
		"status": "healthy",
	})
}

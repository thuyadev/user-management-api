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
// @Param        search    query  string  false  "Search by name"
// @Param        user_id   query  int     false  "Filter by user ID"
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
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	logs, total, err := ctrl.logService.List(c.Request.Context(), uint(userID), page, perPage, search)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, logs, total, page, perPage)
}

// EventStats godoc
// @Summary      Event frequency stats (pie chart)
// @Description  Count of each activity event in the last N days. Requires logs.view permission.
// @Tags         Logs
// @Produce      json
// @Param        days     query  int  false  "Number of days to include (including today)"  default(30)
// @Param        user_id  query  int  false  "Filter by user ID"
// @Success      200  {object}  models.SwaggerLogEventStatsResponse
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      403  {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/logs/stats/events [get]
func (ctrl *LogController) EventStats(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	days := parseStatsDays(c)

	stats, err := ctrl.logService.EventStats(c.Request.Context(), uint(userID), days)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Event statistics retrieved", stats)
}

// DailyStats godoc
// @Summary      Daily activity stats (bar chart)
// @Description  Activity count per day for the last N days (includes quiet days with count 0). Requires logs.view permission.
// @Tags         Logs
// @Produce      json
// @Param        days     query  int  false  "Number of days to include (including today)"  default(30)
// @Param        user_id  query  int  false  "Filter by user ID"
// @Success      200  {object}  models.SwaggerLogDailyStatsResponse
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      403  {object}  models.SwaggerErrorResponse
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Router       /api/v1/admin/logs/stats/daily [get]
func (ctrl *LogController) DailyStats(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	days := parseStatsDays(c)

	stats, err := ctrl.logService.DailyStats(c.Request.Context(), uint(userID), days)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Daily activity statistics retrieved", stats)
}

func parseStatsDays(c *gin.Context) int {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days < 1 {
		days = 1
	}
	if days > 365 {
		days = 365
	}
	return days
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

package models

// Swagger envelope types for OpenAPI documentation.

type SwaggerErrorResponse struct {
	Success bool        `json:"success" example:"false"`
	Message string      `json:"message" example:"Unauthenticated"`
	Errors  interface{} `json:"errors,omitempty"`
}

type SwaggerSuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"OK"`
	Data    interface{} `json:"data,omitempty"`
}

type SwaggerPaginationMeta struct {
	Total      int64 `json:"total" example:"3"`
	Page       int   `json:"page" example:"1"`
	PerPage    int   `json:"per_page" example:"10"`
	TotalPages int   `json:"total_pages" example:"1"`
}

type SwaggerPaginatedUsersResponse struct {
	Success bool           `json:"success" example:"true"`
	Data    []UserResponse `json:"data"`
	Meta    SwaggerPaginationMeta `json:"meta"`
}

type SwaggerPaginatedCategoriesResponse struct {
	Success bool       `json:"success" example:"true"`
	Data    []Category `json:"data"`
	Meta    SwaggerPaginationMeta `json:"meta"`
}

type SwaggerPaginatedProductsResponse struct {
	Success bool      `json:"success" example:"true"`
	Data    []Product `json:"data"`
	Meta    SwaggerPaginationMeta `json:"meta"`
}

type SwaggerPaginatedLogsResponse struct {
	Success bool      `json:"success" example:"true"`
	Data    []UserLog `json:"data"`
	Meta    SwaggerPaginationMeta `json:"meta"`
}

type SwaggerLogEventStatsResponse struct {
	Success bool           `json:"success" example:"true"`
	Message string         `json:"message" example:"Event statistics retrieved"`
	Data    []LogEventStat `json:"data"`
}

type SwaggerLogDailyStatsResponse struct {
	Success bool           `json:"success" example:"true"`
	Message string         `json:"message" example:"Daily activity statistics retrieved"`
	Data    []LogDailyStat `json:"data"`
}

type SwaggerAuthTokenSuccessResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"Login successful"`
	Data    AuthTokenResponse `json:"data"`
}

type SwaggerUserSuccessResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"User retrieved"`
	Data    UserResponse `json:"data"`
}

type SwaggerCategorySuccessResponse struct {
	Success bool     `json:"success" example:"true"`
	Message string   `json:"message" example:"Category retrieved"`
	Data    Category `json:"data"`
}

type SwaggerProductSuccessResponse struct {
	Success bool    `json:"success" example:"true"`
	Message string  `json:"message" example:"Product retrieved"`
	Data    Product `json:"data"`
}

type SwaggerMeResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Profile retrieved"`
	Data    struct {
		UserID      uint     `json:"user_id" example:"1"`
		Email       string   `json:"email" example:"admin@example.com"`
		Role        string   `json:"role" example:"admin"`
		Permissions []string `json:"permissions"`
	} `json:"data"`
}

type SwaggerRolesResponse struct {
	Success bool                `json:"success" example:"true"`
	Message string              `json:"message" example:"Roles and permissions"`
	Data    map[string][]string `json:"data"`
}

type SwaggerHealthResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"API is running"`
	Data    struct {
		Status string `json:"status" example:"healthy"`
	} `json:"data"`
}

type SwaggerAISuggestResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Category name suggested"`
	Data    struct {
		Name string `json:"name" example:"Fitness & Gym Equipment"`
	} `json:"data"`
}

type SwaggerAIDescriptionResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Description generated"`
	Data    struct {
		Description string `json:"description"`
	} `json:"data"`
}

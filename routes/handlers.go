package routes

import "user-management-api/controllers"

type Handlers struct {
	Auth     *controllers.AuthController
	User     *controllers.UserController
	Category *controllers.CategoryController
	Product  *controllers.ProductController
	Log      *controllers.LogController
}

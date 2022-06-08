package routes

import (
	"server/controllers"
	"server/middlewares"
	"server/services"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.POST("/signup", controllers.SignUp())
	router.POST("/signin", controllers.SignIn())
	router.GET("/google/:email", services.HandleGoogleLogin())
	router.GET("/signin-google", services.CallBackFromGoogle())
}

func DriveRoutes(router *gin.Engine) {
	router.POST("/registercar", controllers.CreateCar())
	router.GET("/profilecar", controllers.GetCarProfileByEmail())
	router.POST("/review", controllers.CreateReview())
}

func DriveRoute(router *gin.Engine) {
	UserRoutes(router)
	router.Use(middlewares.DeserializeUser())
	DriveRoutes(router)
	// router.POST("/user", controllers.CreateUser())
	// router.GET("/user/:userId", controllers.GetAUser())
	// router.PUT("/user/:userId", controllers.EditAUser())
	// router.DELETE("/user/:userId", controllers.DeleteAUser())
	// router.GET("/users", controllers.GetAllUsers())
	// router.POST("/register", contro)
}

package routes

import (
	"net/http"
	"taskmanagement/controllers"
	"taskmanagement/middlewares"

	websocket "taskmanagement/socket"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"}, // Frontend URL
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true, // Allow cookies, authorization headers
	// }))

	authRoutes := r.Group("/auth")
	authRoutes.POST("/register", controllers.Register)
	authRoutes.POST("/login", controllers.Login)

	r.GET("/ws", websocket.HandleWebSocket)

	protectedRoutes := r.Group("/api")
	protectedRoutes.Use(middlewares.JWTAuthMiddleware())

	protectedRoutes.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "You are authorized!"})
	})

	taskRoutes := protectedRoutes.Group("/tasks")
	{
		taskRoutes.POST("/breakdown", controllers.BreakDown)
		taskRoutes.POST("/suggest", controllers.SuggestTask)
		taskRoutes.POST("/priotize", controllers.ProtizeTask)

	}

	userRoutes := protectedRoutes.Group("/users")
	{
		userRoutes.GET("/", controllers.GetAllUsers)
		userRoutes.POST("/user", controllers.GetUser)
	}

	return r
}

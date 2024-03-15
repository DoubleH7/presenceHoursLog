package main

import (
	"fmt"
	"os"

	"github.com/DoubleH7/presenceHoursLog/database"
	"github.com/DoubleH7/presenceHoursLog/webService"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	fmt.Println("initializing")

	godotenv.Load(".env")
	port := os.Getenv("PORT")

	fmt.Print("Connecting to database ...")

	client, err := database.ConnectDB()
	defer database.DisconnectDB(client)

	if err != nil {
		fmt.Println("failed")
		panic(err)
	}

	fmt.Println("\nStarting server...")

	e := echo.New()

	adminGroup := e.Group("/admin")

	// adding basic authentication to various endpoints
	adminGroup.Use(middleware.BasicAuth(webService.UserpassCheck))

	// Setting up admin logs
	fmt.Print("\tadmin log setup...")
	file, err := os.OpenFile("./admin_access.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("failed")
		panic(err)
	}
	fmt.Println("")

	adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}]   ${Status} from ${remote_ip}  ${method}${path}  ${latency_human}` + "\n",
		Output: file,
	}))

	//hooking up the server check handler
	e.GET("/", webService.ServerAlive(client))

	// hooking up the admin handlers
	adminGroup.GET("/users/all", webService.GetUsers(client))
	adminGroup.GET("/users/id/:id", webService.GetUserbyid(client))
	adminGroup.GET("/users/name/:name", webService.GetUserbyname(client))
	adminGroup.POST("/users", webService.CreateUser(client))

	// hooking up the user handlers
	e.POST("/start/:id", webService.CreateSitting(client))
	e.POST("/stop/:id", webService.StopSitting(client))

	//starting server
	fmt.Println("\tinitializing...")
	err = e.Start(fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println("Failed")
		panic(err)
	}
}

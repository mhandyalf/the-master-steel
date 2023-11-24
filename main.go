package main

import (
	"fmt"
	"os"
	"the-master-steel/database"
	"the-master-steel/handlers"
	"the-master-steel/middleware"

	"github.com/labstack/echo/v4"
	_ "github.com/joho/godotenv/autoload"

)

func main() {
	e := echo.New()

	authhandler := handlers.NewAuth(database.InitDB())
	e.POST("/register", authhandler.Register)
	e.POST("/login", authhandler.Login)
	e.GET("/employee", authhandler.GetEmployeeInfo, middleware.JWTAuth)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))

}

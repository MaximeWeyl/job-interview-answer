package main

//go:generate swag init --parseDependency

// @title Falcon Backend API
// @version 1.0
// @description This is the backend for the misteryemployer challenge

// @contact.name Maxime Weyl
// @contact.email weyl.maxime@gmail.com

// @BasePath /api/v1

// giveThemTheOddsInput godoc
// @Summary List accounts
// @Description get accounts
// @Accept  json
// @Produce  json
// @Param q query string false "name search by q"
// @Success 200 {array} model.Account
// @Header 200 {string} Token "qwerty"
// @Failure 400,404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Failure default {object} httputil.DefaultError
// @Router /accounts [get]

import (
	"fmt"
	"os"

	_ "github.com/MaximeWeyl/misteryemployer-what-are-the-odds/backend/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cmd := cobra.Command{
		Use:   "falcon millenium_falcon_file",
		Short: "Runs a backend ready to give the odds !",
		Long: `This is the falcon's onboard computer, ready
to serve information directly in C3PO's ship over the air.
This listens on the port 8080 by default (can be overridden
with the PORT environment variable).

,---,---,---,---,---,---,---,---,---,---,---,---,---,-------,
|1/2| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 0 | + | ' | <-    |
|---'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-----|
| ->| | Q | W | E | R | T | Y | U | I | O | P | ] | ^ |     |
|-----',--',--',--',--',--',--',--',--',--',--',--',--'|    |
| Caps | A | S | D | F | G | H | J | K | L | \ | [ | * |    |
|----,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'---'----|
|    | < | Z | X | C | V | B | N | M | , | . | - |          |
|----'-,-',--'--,'---'---'---'---'---'---'-,-'---',--,------|
| ctrl |  | alt |                          |altgr |  | ctrl |
'------'  '-----'--------------------------'------'  '------'

`,
		Args: cobra.ExactArgs(1),
		RunE: run,
	}
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	err := initFromCobra(cmd, args)
	if err != nil {
		return err
	}

	err = runServer()
	if err != nil {
		return err
	}

	return nil
}

func runServer() error {
	// Dev or Production mode based on environment variables
	_, devMode := os.LookupEnv("DEBUG_MODE")
	if devMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	apiV1 := r.Group("/api/v1")

	// CORS : this should always be done before adding any route
	fmt.Println("Allowing all CORS origins")
	config := cors.Default()
	r.Use(config)     // The middleware must be used both on the default handler
	apiV1.Use(config) // and the group

	// ROUTES
	// API
	apiV1.POST("/give-me-the-odds", GiveThemTheOdds)

	// use ginSwagger middleware to serve the API docs
	apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Listen and serve on port the port as defined by the
	// environment variable, with 8080 as a default
	portString := os.Getenv("PORT")
	if portString == "" {
		portString = "8080"
	}
	err := r.Run(fmt.Sprintf("0.0.0.0:%s", portString))
	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"fmt"
	"net/http"

	"code-executor/internal/api"
	"code-executor/internal/storage"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Code Executor API
// @version 1.0
// @description This is a code execution service that allows submitting programs for execution and retrieving results asynchronously.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	storageInstance := storage.NewInMemoryStorage()
	api.Init(storageInstance)

	http.HandleFunc("/task", api.LoadTask)
	http.HandleFunc("/status/", api.CheckTaskStatus)
	http.HandleFunc("/result/", api.GetResult)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("an error occured while trying to start the server: ", err)
	}
}

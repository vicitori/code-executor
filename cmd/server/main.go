package main

import (
	"fmt"
	"net/http"

	_ "code-executor/docs"
	"code-executor/internal/api"
	"code-executor/internal/storage"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Code Executor API
// @version 1.0
// @host localhost:8080
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

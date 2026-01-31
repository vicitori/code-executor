package main

import (
	"fmt"
	"net/http"

	"code-executor/internal/api"
)

func main() {
	http.HandleFunc("/task", api.LoadTask)
	http.HandleFunc("/status/", api.CheckTaskStatus)
	http.HandleFunc("/result/", api.GetResult)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("an error occured while trying to start the server: ", err)
	}
}

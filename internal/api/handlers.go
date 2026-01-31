package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"task-service/internal/domain"
	"task-service/internal/storage"
)

var tS storage.TaskStorage

func Init(storage storage.TaskStorage) {
	tS = storage
}

func LoadTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, `{"error": "app supports only Content-Type application/json"}`, http.StatusUnsupportedMediaType)
		return
	}

	defer r.Body.Close()
	var req domain.Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error": "invalid JSON: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	if req.Program == "" {
		http.Error(w, `{"error": "program field is required"}`, http.StatusBadRequest)
		return
	}

	if req.Compiler == "" {
		http.Error(w, `{"error": "compiler field is required"}`, http.StatusBadRequest)
		return
	}

	taskId, err := tS.CreateTask(req.Program, req.Compiler)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	go processTask(taskId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	var resp domain.IdResponse
	resp.Id = taskId
	json.NewEncoder(w).Encode(&resp)
}

func CheckTaskStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	taskId := strings.TrimPrefix(r.URL.Path, "/status/")
	if taskId == "" {
		http.Error(w, `{"error": "task id is required"}`, http.StatusBadRequest)
		return
	}
	task, err := tS.GetTask(taskId)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var resp domain.StatusResponse
	resp.Status = task.Status
	json.NewEncoder(w).Encode(&resp)
}

func GetResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	taskId := strings.TrimPrefix(r.URL.Path, "/result/")
	if taskId == "" {
		http.Error(w, `{"error": "task id is required"}`, http.StatusBadRequest)
		return
	}
	task, err := tS.GetTask(taskId)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	if task.Status != domain.Ready {
		http.Error(w, `{"error": "task is not ready yet"}`, http.StatusTooEarly)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var resp domain.ResultResponse
	resp.Result = task.Result
	json.NewEncoder(w).Encode(&resp)
}

func processTask(taskID string) {
	tS.UpdateStatus(taskID, domain.InProgress)
	time.Sleep(10 * time.Second)
	result := fmt.Sprintf("Processed task %s at %s.",
		taskID,
		time.Now().Format(time.RFC3339))
	tS.SetResult(taskID, result)
}

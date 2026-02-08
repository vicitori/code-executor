package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"code-executor/internal/domain"
	"code-executor/internal/storage"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, message string, status int) {
	writeJSON(w, status, domain.ErrorResponse{Error: message})
}

var tS storage.TaskStorage

func Init(storage storage.TaskStorage) {
	tS = storage
}

// LoadTask godoc
// @Accept json
// @Produce json
// @Param request body domain.Request true "Task data"
// @Success 201 {object} domain.IdResponse
// @Router /task [post]
func LoadTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		writeError(w, "app supports only Content-Type application/json", http.StatusUnsupportedMediaType)
		return
	}

	defer r.Body.Close()
	var req domain.Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Program == "" {
		writeError(w, "program field is required", http.StatusBadRequest)
		return
	}

	if req.Compiler == "" {
		writeError(w, "compiler field is required", http.StatusBadRequest)
		return
	}

	taskId, err := tS.CreateTask(req.Program, req.Compiler)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go processTask(taskId)

	writeJSON(w, http.StatusCreated, domain.IdResponse{Id: taskId})
}

// CheckTaskStatus godoc
// @Produce json
// @Param task_id path string true "Task ID"
// @Success 200 {object} domain.StatusResponse
// @Router /status/{task_id} [get]
func CheckTaskStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	taskId := strings.TrimPrefix(r.URL.Path, "/status/")
	if taskId == "" {
		writeError(w, "task id is required", http.StatusBadRequest)
		return
	}
	task, err := tS.GetTask(taskId)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, domain.StatusResponse{Status: task.Status})
}

// GetResult godoc
// @Produce json
// @Param task_id path string true "Task ID"
// @Success 200 {object} domain.ResultResponse
// @Router /result/{task_id} [get]
func GetResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	taskId := strings.TrimPrefix(r.URL.Path, "/result/")
	if taskId == "" {
		writeError(w, "task id is required", http.StatusBadRequest)
		return
	}
	task, err := tS.GetTask(taskId)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	if task.Status != domain.Ready {
		writeError(w, "task is not ready yet", http.StatusTooEarly)
		return
	}

	writeJSON(w, http.StatusOK, domain.ResultResponse{Result: task.Result})
}

func processTask(taskID string) {
	tS.UpdateStatus(taskID, domain.InProgress)
	time.Sleep(1 * time.Second)
	result := fmt.Sprintf("processed task %s at %s",
		taskID,
		time.Now().Format(time.RFC3339))
	tS.SetResult(taskID, result)
}

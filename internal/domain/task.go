package domain

type TaskStatus string

const (
	Created    TaskStatus = "created"
	InProgress TaskStatus = "in_progress"
	Ready      TaskStatus = "ready"
	Failed     TaskStatus = "failed"
)

type Task struct {
	Id       string     `json:"id"`
	Program  string     `json:"program"`
	Compiler string     `json:"compiler"`
	Status   TaskStatus `json:"status"`
	Result   string     `json:"result"`
}

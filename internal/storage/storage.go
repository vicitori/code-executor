package storage

import "code-executor/internal/domain"

type TaskStorage interface {
	CreateTask(userProgram, compilerName string) (string, error)
	UpdateStatus(uuid string, status domain.TaskStatus) error
	GetResult(uuid string) (string, error)
	GetTask(uuid string) (*domain.Task, error)
	SetResult(uuid, result string) error
}

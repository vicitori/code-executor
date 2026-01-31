package storage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"code-executor/internal/domain"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type InMemoryStorage struct {
	mutex sync.RWMutex
	tasks map[string]*domain.Task
	genId func() string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		tasks: make(map[string]*domain.Task),
		genId: func() string {
			return uuid.New().String()
		},
	}
}

func (st *InMemoryStorage) CreateTask(program, compiler string) (string, error) {
	id := st.genId()
	task := &domain.Task{
		Id:       id,
		Program:  program,
		Compiler: compiler,
		Status:   domain.Created,
		Result:   "",
	}
	st.mutex.Lock()
	st.tasks[id] = task
	st.mutex.Unlock()
	return id, nil
}

func (st *InMemoryStorage) GetTask(uuid string) (*domain.Task, error) {
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	task, exists := st.tasks[uuid]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrTaskNotFound, uuid)
	}
	return task, nil
}

func (st *InMemoryStorage) UpdateStatus(uuid string, status domain.TaskStatus) error {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	task, exists := st.tasks[uuid]
	if !exists {
		return fmt.Errorf("%w: %s", ErrTaskNotFound, uuid)
	}
	task.Status = status
	return nil
}

func (st *InMemoryStorage) GetResult(uuid string) (string, error) {
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	task, exists := st.tasks[uuid]
	if !exists {
		return "", fmt.Errorf("%w: %s", ErrTaskNotFound, uuid)
	}
	res := task.Result
	return res, nil
}

func (st *InMemoryStorage) SetResult(uuid, result string) error {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	task, exists := st.tasks[uuid]
	if !exists {
		return fmt.Errorf("%w: %s", ErrTaskNotFound, uuid)
	}
	task.Result = result
	task.Status = domain.Ready
	return nil
}

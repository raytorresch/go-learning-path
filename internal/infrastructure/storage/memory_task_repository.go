package storage

import (
	"fmt"
	"sync"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/repositories"
)

type MemoryTaskRespository struct {
	tasks  map[int]*entities.Task
	mutex  sync.RWMutex
	nextID int
}

var _ repositories.TaskRepository = (*MemoryTaskRespository)(nil)

func NewMemoryTaskRepository() *MemoryTaskRespository {
	return &MemoryTaskRespository{
		tasks:  make(map[int]*entities.Task),
		nextID: 1,
	}
}

func (r *MemoryTaskRespository) Save(task *entities.Task) (*entities.Task, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	task.ID = r.nextID
	r.tasks[task.ID] = task
	r.nextID++
	return task, nil
}

func (r *MemoryTaskRespository) FindByUserID(userID int) ([]*entities.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	var tasks []*entities.Task
	for _, task := range r.tasks {
		if task.UserID == userID {
			tasks = append(tasks, task)
		}
	}
	if len(tasks) == 0 {
		return nil, fmt.Errorf("tarea no encontrada")
	}
	return tasks, nil
}

package repositories

import (
	"user-management/internal/domain/entities"
)

type TaskRepository interface {
	Save(task *entities.Task) (*entities.Task, error)
	FindByUserID(userID int) ([]*entities.Task, error)
}

// type TaskRepository struct {
// 	tasks  map[int]*entities.Task
// 	mutex  sync.RWMutex
// 	nextID int
// }

// func NewTaskRepository() *TaskRepository {
// 	return &TaskRepository{
// 		tasks:  make(map[int]*entities.Task),
// 		nextID: 1,
// 	}
// }

// func (r *TaskRepository) Save(task *entities.Task) (*entities.Task, error) {
// 	r.mutex.Lock()
// 	defer r.mutex.Unlock()

// 	task.ID = r.nextID
// 	r.tasks[task.ID] = task
// 	r.nextID++
// 	return task, nil
// }

// func (r *TaskRepository) FindByID(id int) (*entities.Task, error) {
// 	r.mutex.RLock()
// 	defer r.mutex.RUnlock()

// 	task, exists := r.tasks[id]
// 	if !exists {
// 		return nil, fmt.Errorf("tarea no encontrada")
// 	}
// 	return task, nil
// }

// func (r *TaskRepository) FindAll() []*entities.Task {
// 	r.mutex.RLock()
// 	defer r.mutex.RUnlock()

// 	tasks := make([]*entities.Task, 0, len(r.tasks))
// 	for _, task := range r.tasks {
// 		tasks = append(tasks, task)
// 	}
// 	return tasks
// }

// func (r *TaskRepository) Update(task *entities.Task) (*entities.Task, error) {
// 	r.mutex.Lock()
// 	defer r.mutex.Unlock()

// 	if _, exists := r.tasks[task.ID]; !exists {
// 		return nil, fmt.Errorf("tarea no encontrada")
// 	}
// 	r.tasks[task.ID] = task
// 	return task, nil
// }

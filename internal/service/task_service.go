package service

import (
	"user-management/internal/models"
	"user-management/internal/processors"
	"user-management/internal/storage"
)

type TaskService struct {
	repo      *storage.TaskRepository
	processor *processors.TaskProcessor
}

func NewTaskService(repo *storage.TaskRepository, processor *processors.TaskProcessor) *TaskService {
	service := &TaskService{
		repo:      repo,
		processor: processor,
	}
	service.registerHandlers()
	return service
}

func (s *TaskService) registerHandlers() {
	// Registrando handlers como closures
	s.processor.RegisterHandler("email_report", func(task *models.Task) error {
		// Simular procesamiento
		task.Data = "reporte_generado"
		return nil
	})

	s.processor.RegisterHandler("data_cleanup", func(task *models.Task) error {
		// Simular limpieza de datos
		task.Data = "datos_limpiados"
		return nil
	})
}

// CreateAndProcess demuestra control de flujo
func (s *TaskService) CreateAndProcess(name string, userID int, priority int) (*models.Task, error) {
	task := &models.Task{
		UserID:   userID,
		Name:     name,
		Priority: priority,
		Status:   models.StatusPending,
	}

	// Guardar tarea
	savedTask, err := s.repo.Save(task)
	if err != nil {
		return nil, err
	}

	// Procesar tarea
	err = s.processor.ProcessTask(savedTask)
	if err != nil {
		// Intentar nuevamente con prioridad aumentada
		savedTask.Priority++
		_ = s.processor.ProcessTask(savedTask)
	}

	return s.repo.Update(savedTask)
}

// BatchProcess demuestra funciones variádicas
func (s *TaskService) BatchProcess(tasks ...*models.Task) (completed int, failed int) {
	return s.processor.ProcessBatch(tasks...)
}

// GetHighPriorityTasks muestra filtros
func (s *TaskService) GetHighPriorityTasks(minPriority int) []*models.Task {
	allTasks := s.repo.FindAll()
	filter := processors.CreatePriorityFilter(minPriority)
	return s.processor.FilterTasks(allTasks, filter)
}

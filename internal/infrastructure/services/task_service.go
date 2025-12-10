package services

import (
	"user-management/internal/domain/repositories"
	"user-management/internal/processors"
)

type TaskService struct {
	repo      *repositories.TaskRepository
	processor *processors.TaskProcessor
}

func NewTaskService(repo *repositories.TaskRepository, processor *processors.TaskProcessor) *TaskService {
	service := &TaskService{
		repo:      repo,
		processor: processor,
	}
	return service
}

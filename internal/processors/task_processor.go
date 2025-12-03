package processors

import (
	"fmt"
	"log"
	"time"
	"user-management/internal/models"
)

// Type aliases para funciones (first-class citizens)
type TaskHandler func(task *models.Task) error
type TaskFilter func(task *models.Task) bool

// TaskProcessor demuestra control de flujo avanzado
type TaskProcessor struct {
	handlers map[string]TaskHandler
}

func NewTaskProcessor() *TaskProcessor {
	return &TaskProcessor{
		handlers: make(map[string]TaskHandler),
	}
}

// RegisterHandler muestra funciones como valores
func (p *TaskProcessor) RegisterHandler(taskType string, handler TaskHandler) {
	p.handlers[taskType] = handler
}

// ProcessBatch demuestra funciones variádicas
func (p *TaskProcessor) ProcessBatch(tasks ...*models.Task) (int, int) {
	completed := 0
	failed := 0

	for i, task := range tasks {
		// Defer para logging (ejecuta al final de cada iteración)
		defer func(idx int) {
			log.Printf("Procesada tarea %d: %s", idx+1, tasks[idx].Name)
		}(i)

		err := p.ProcessTask(task)
		if err != nil {
			failed++
			log.Printf("Error procesando tarea %s: %v", task.Name, err)
		} else {
			completed++
		}
	}

	return completed, failed
}

// ProcessTask con defer, panic y recover
func (p *TaskProcessor) ProcessTask(task *models.Task) (err error) {
	// Defer para manejo de panic
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recuperado: %v", r)
			task.Status = models.StatusFailed
		}
	}()

	// Validación con switch avanzado
	switch {
	case task.Name == "":
		return fmt.Errorf("nombre de tarea requerido")
	case task.Priority < 1 || task.Priority > 5:
		task.Priority = 3 // Valor por defecto usando fallthrough
		fallthrough
	default:
		task.Status = models.StatusProcessing
	}

	// Closure para tiempo de procesamiento
	processTime := func() string {
		start := time.Now()
		defer func() {
			elapsed := time.Since(start)
			log.Printf("Tiempo de procesamiento: %v", elapsed)
		}()
		return "procesado"
	}

	_ = processTime() // Ejecutar closure

	// Handler específico por tipo (simulado)
	handler, exists := p.handlers[task.Name]
	if exists {
		return handler(task)
	}

	// Handler por defecto
	task.Status = models.StatusCompleted
	return nil
}

// FilterTasks muestra funciones como parámetros
func (p *TaskProcessor) FilterTasks(tasks []*models.Task, filter TaskFilter) []*models.Task {
	var result []*models.Task
	for _, task := range tasks {
		if filter(task) {
			result = append(result, task)
		}
	}
	return result
}

// Crear filtros usando closures
func CreatePriorityFilter(minPriority int) TaskFilter {
	return func(task *models.Task) bool {
		return task.Priority >= minPriority
	}
}

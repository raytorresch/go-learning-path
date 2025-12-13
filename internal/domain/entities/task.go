package entities

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
)

type Task struct {
	ID       int        `json:"id"`
	UserID   int        `json:"user_id"`
	Name     string     `json:"name"`
	Status   TaskStatus `json:"status"`
	Priority int        `json:"priority"` // 1-5, donde 5 es más alto
	Data     any        `json:"data"`     // any para datos flexibles
}

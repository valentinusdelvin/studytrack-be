package enums

type TaskDeadlineStatus string

const (
	TaskDeadlineStatusNormal  TaskDeadlineStatus = "normal"
	TaskDeadlineStatusDueSoon TaskDeadlineStatus = "due_soon"
	TaskDeadlineStatusOverdue TaskDeadlineStatus = "overdue"
)

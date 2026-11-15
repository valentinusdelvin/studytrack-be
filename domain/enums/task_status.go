package enums

import (
	"database/sql/driver"
	"fmt"
)

type TaskStatus string

const (
	TaskStatusDone    TaskStatus = "selesai"
	TaskStatusNotDone TaskStatus = "belum"
)

func (t *TaskStatus) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*t = TaskStatus(v)
	case string:
		*t = TaskStatus(v)
	case nil:
		return nil
	default:
		return fmt.Errorf("unsupported TaskStatus type: %T", value)
	}
	return nil
}

func (t TaskStatus) Value() (driver.Value, error) {
	return string(t), nil
}

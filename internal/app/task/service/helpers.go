package service

import (
	"slices"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/valentinusdelvin/studytrack-be/domain/dto"
	"github.com/valentinusdelvin/studytrack-be/domain/enums"
	"github.com/valentinusdelvin/studytrack-be/domain/errx"
	"github.com/valentinusdelvin/studytrack-be/internal/app/task/repository"
)

const taskDueSoonThreshold = 24 * time.Hour

var taskDeadlineLocation = func() *time.Location {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}
	return location
}()

func formatTaskListItems(rows []repository.TaskWithCourse, now time.Time) []dto.TaskListItem {
	type taskListEntry struct {
		item     dto.TaskListItem
		deadline time.Time
	}

	entries := make([]taskListEntry, 0, len(rows))
	for _, row := range rows {
		deadline := row.Deadline
		deadlineStatus := getTaskDeadlineStatus(row.Status, deadline, now)
		item := dto.FormatToTaskListItemResponse(row.Task, row.CourseName, deadlineStatus)
		entries = append(entries, taskListEntry{item: item, deadline: deadline})
	}

	dueSoon := string(enums.TaskDeadlineStatusDueSoon)
	slices.SortStableFunc(entries, func(first, second taskListEntry) int {
		firstDueSoon := first.item.DeadlineStatus == dueSoon
		secondDueSoon := second.item.DeadlineStatus == dueSoon
		switch {
		case firstDueSoon && secondDueSoon:
			return first.deadline.Compare(second.deadline)
		case firstDueSoon:
			return -1
		case secondDueSoon:
			return 1
		default:
			return 0
		}
	})

	out := make([]dto.TaskListItem, 0, len(entries))
	for _, entry := range entries {
		out = append(out, entry.item)
	}
	return out
}

func parseTaskDeadline(value string) (time.Time, error) {
	for _, layout := range []string{"2006-01-02T15:04", "2006-01-02T15:04:05", time.RFC3339, time.RFC3339Nano} {
		if parsed, err := time.ParseInLocation(layout, value, taskDeadlineLocation); err == nil {
			return parsed.In(taskDeadlineLocation), nil
		}
	}
	return time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid deadline")
}

func getTaskDeadlineStatus(status enums.TaskStatus, deadline, now time.Time) enums.TaskDeadlineStatus {
	if status == enums.TaskStatusDone {
		return enums.TaskDeadlineStatusNormal
	}

	remaining := deadline.Sub(now)
	switch {
	case remaining < 0:
		return enums.TaskDeadlineStatusOverdue
	case remaining < taskDueSoonThreshold:
		return enums.TaskDeadlineStatusDueSoon
	default:
		return enums.TaskDeadlineStatusNormal
	}
}

func dbError(err error, location string) error {
	return errx.ErrInternalServer.WithError(err).WithLocation(location)
}

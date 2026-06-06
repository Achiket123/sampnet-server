package task_comments

import (
	"context"
	"fmt"
	"log"
	"server/internal/domain/notifications"
	"server/internal/domain/task_comments"
	"server/internal/domain/tasks"
	"time"
)

type service struct {
	repo           task_comments.Repository
	taskRepo       tasks.Repository
	notificationUc notifications.UseCase
}

func NewService(repo task_comments.Repository, taskRepo tasks.Repository, notificationUc notifications.UseCase) task_comments.UseCase {
	return &service{
		repo:           repo,
		taskRepo:       taskRepo,
		notificationUc: notificationUc,
	}
}

func (s *service) AddComment(ctx context.Context, taskID uint, userID uint, content string) (*task_comments.TaskComment, error) {
	comment := &task_comments.TaskComment{
		TaskID:    taskID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, err
	}

	go func() {
		task, err := s.taskRepo.GetByID(context.Background(), taskID)
		if err != nil || task == nil {
			log.Printf("failed to fetch task for comment notification: %v", err)
			return
		}

		recipients := make(map[uint]bool)
		recipients[task.AssignedBy] = true
		recipients[task.AssignedTo] = true

		for recipientID := range recipients {
			if recipientID == userID || recipientID == 0 {
				continue
			}

			title := "New comment on task"
			message := fmt.Sprintf("A new comment was added to task: %s", task.Title)
			link := fmt.Sprintf("/task-detail/%d", task.ID)
			_ = s.notificationUc.CreateNotification(context.Background(), recipientID, task.OrganisationID, title, message, "task_comment", link)
		}
	}()

	return comment, nil
}

func (s *service) GetComments(ctx context.Context, taskID uint) ([]task_comments.TaskComment, error) {
	return s.repo.GetByTaskID(ctx, taskID)
}

func (s *service) DeleteComment(ctx context.Context, commentID uint, userID uint) error {
	comment, err := s.repo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return task_comments.ErrCommentNotFound
	}

	if comment.UserID != userID {
		return task_comments.ErrNotCommentOwner
	}

	return s.repo.Delete(ctx, commentID)
}

package task_comments

import (
	"context"
	"server/internal/domain/task_comments"
	"time"
)

type service struct {
	repo task_comments.Repository
}

func NewService(repo task_comments.Repository) task_comments.UseCase {
	return &service{repo: repo}
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

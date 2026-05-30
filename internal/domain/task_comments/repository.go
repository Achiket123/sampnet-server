package task_comments

import (
	"context"
	"errors"
)

var (
	ErrNotCommentOwner = errors.New("not the owner of the comment")
	ErrCommentNotFound = errors.New("comment not found")
)

type Repository interface {
	Create(ctx context.Context, comment *TaskComment) error
	GetByTaskID(ctx context.Context, taskID uint) ([]TaskComment, error)
	Delete(ctx context.Context, commentID uint) error
	GetByID(ctx context.Context, commentID uint) (*TaskComment, error)
}

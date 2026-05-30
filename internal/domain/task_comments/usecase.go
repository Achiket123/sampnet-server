package task_comments

import (
	"context"
)

type UseCase interface {
	AddComment(ctx context.Context, taskID uint, userID uint, content string) (*TaskComment, error)
	GetComments(ctx context.Context, taskID uint) ([]TaskComment, error)
	DeleteComment(ctx context.Context, commentID uint, userID uint) error
}

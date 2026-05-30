package task_attachments

import (
	"context"
	"errors"
)

var (
	ErrNotAttachmentOwner = errors.New("not the owner of the attachment")
	ErrAttachmentNotFound = errors.New("attachment not found")
)

type Repository interface {
	Create(ctx context.Context, attachment *TaskAttachment) error
	GetByTaskID(ctx context.Context, taskID uint) ([]TaskAttachment, error)
	Delete(ctx context.Context, attachmentID uint) error
	GetByID(ctx context.Context, attachmentID uint) (*TaskAttachment, error)
}

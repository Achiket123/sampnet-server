package task_attachments

import (
	"context"
)

type UseCase interface {
	AttachFile(ctx context.Context, taskID uint, fileID uint, uploadedBy uint, fileName string) (*TaskAttachment, error)
	GetAttachments(ctx context.Context, taskID uint) ([]TaskAttachment, error)
	RemoveAttachment(ctx context.Context, attachmentID uint, requestingUserID uint) error
}

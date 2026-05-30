package task_attachments

import (
	"context"
	"server/internal/domain/task_attachments"
	"time"
)

type service struct {
	repo task_attachments.Repository
}

func NewService(repo task_attachments.Repository) task_attachments.UseCase {
	return &service{repo: repo}
}

func (s *service) AttachFile(ctx context.Context, taskID uint, fileID uint, uploadedBy uint, fileName string) (*task_attachments.TaskAttachment, error) {
	attachment := &task_attachments.TaskAttachment{
		TaskID:     taskID,
		FileID:     fileID,
		UploadedBy: uploadedBy,
		FileName:   fileName,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, attachment); err != nil {
		return nil, err
	}

	return attachment, nil
}

func (s *service) GetAttachments(ctx context.Context, taskID uint) ([]task_attachments.TaskAttachment, error) {
	return s.repo.GetByTaskID(ctx, taskID)
}

func (s *service) RemoveAttachment(ctx context.Context, attachmentID uint, requestingUserID uint) error {
	attachment, err := s.repo.GetByID(ctx, attachmentID)
	if err != nil {
		return err
	}
	if attachment == nil {
		return task_attachments.ErrAttachmentNotFound
	}

	if attachment.UploadedBy != requestingUserID {
		return task_attachments.ErrNotAttachmentOwner
	}

	return s.repo.Delete(ctx, attachmentID)
}

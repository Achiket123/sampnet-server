package task_attachments

import "time"

type TaskAttachment struct {
	ID                 uint      `json:"id"`
	TaskID             uint      `json:"task_id"`
	FileID             uint      `json:"file_id"`
	UploadedBy         uint      `json:"uploaded_by"`
	FileName           string    `json:"file_name"`
	CreatedAt          time.Time `json:"created_at"`
	UploaderFirstName string    `json:"uploader_first_name"`
	UploaderLastName  string    `json:"uploader_last_name"`
}

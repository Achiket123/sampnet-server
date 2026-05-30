package task_comments

import "time"

type TaskComment struct {
	ID            uint      `json:"id"`
	TaskID        uint      `json:"task_id"`
	UserID        uint      `json:"user_id"`
	UserFirstName string    `json:"user_first_name"`
	UserLastName  string    `json:"user_last_name"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"created_at"`
}

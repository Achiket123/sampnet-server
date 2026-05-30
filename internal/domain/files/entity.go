package files

type File struct {
	ID       uint   `json:"id"`
	FileName string `json:"file_name"`
	Data     []byte `json:"-"`
	FileType string `json:"file_type"`
	FileSize int64  `json:"file_size"`
}

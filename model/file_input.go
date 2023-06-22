package model

import "fmt"

type FileInput struct {
	Name        string
	OwnerID     int
	Content     []byte
	ContentType string
}

func (fileInput FileInput) GetPath() string {
	return fmt.Sprintf("%d/%s", fileInput.OwnerID, fileInput.Name)
}

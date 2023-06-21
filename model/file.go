package model

import (
	"cloud.google.com/go/storage"
	"time"
)

type File struct {
	Name        string
	Size        int64
	ContentType string
	CreatedOn   time.Time
	ModifiedOn  time.Time
	OwnerID     int
}

func FromGCSObject(attrs *storage.ObjectAttrs) File {
	return File{
		Name:        attrs.Name,
		Size:        attrs.Size,
		ContentType: attrs.ContentType,
		CreatedOn:   attrs.Created,
		ModifiedOn:  attrs.Updated,
	}
}

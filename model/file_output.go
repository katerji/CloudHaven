package model

type FileOutput struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	CreatedOn   string `json:"created_on"`
	ModifiedOn  string `json:"modified_on"`
	OwnerID     int    `json:"owner_id"`
}

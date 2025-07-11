package response

type DocumentUploadData struct {
	JSON map[string]interface{} `json:"json,omitempty"`
	File string                 `json:"file"`
}

type DocumentUploadResponse struct {
	Data DocumentUploadData `json:"data"`
}

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/DENFNC/web-test/internal/transport/dto/request"
)

func ParseMeta(r *http.Request) (*request.DocumentMetaRequest, error) {
	metaStr := r.FormValue("meta")
	if metaStr == "" {
		return nil, errors.New("meta is required")
	}
	var meta request.DocumentMetaRequest
	if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
		return nil, errors.New("invalid meta json")
	}
	return &meta, nil
}

func SaveUploadedFile(r *http.Request) (*multipart.FileHeader, error) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("missing file")
	}
	defer file.Close()

	if err := os.MkdirAll("uploads", 0755); err != nil {
		return nil, fmt.Errorf("cannot create upload dir")
	}

	dst, err := os.Create(filepath.Join("uploads", fileHeader.Filename))
	if err != nil {
		return nil, fmt.Errorf("cannot save file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("cannot write file")
	}

	return fileHeader, nil
}

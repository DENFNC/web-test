package handler

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/DENFNC/web-test/internal/domain"
	"github.com/DENFNC/web-test/internal/service"
	"github.com/DENFNC/web-test/internal/transport/dto/response"
	"github.com/DENFNC/web-test/internal/utils"
)

type DocumentService interface{}

type DocumentHandler struct {
	*slog.Logger
	Service *service.DocumentService
}

func NewDocumentHandler(log *slog.Logger, mux *http.ServeMux, docService *service.DocumentService) {
	handler := &DocumentHandler{
		Logger:  log,
		Service: docService,
	}

	mux.HandleFunc("POST /api/docs", handler.createDocumentHandler)
	mux.HandleFunc("GET /api/docs", handler.getDocumentsHandler)
	mux.HandleFunc("GET /api/docs/{id}", handler.getDocumentHandler)
	mux.HandleFunc("DELETE /api/docs/{id}", handler.deleteDocumentHandler)
}

func (api *DocumentHandler) createDocumentHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "Error multipart data")
		return
	}

	meta, err := utils.ParseMeta(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	ownerID, err := api.Service.ValidateToken(r.Context(), meta.Token)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer file.Close()

	originalName := fileHeader.Filename
	docID, uuidFileName, err := api.Service.SaveDocument(r.Context(), *meta, ownerID, originalName)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "cannot save document")
		return
	}

	uploadPath := filepath.Join("data/uploads", uuidFileName)
	out, err := os.Create(uploadPath)
	if err != nil {
		fmt.Println(err)
		response.Error(w, http.StatusInternalServerError, "cannot create file")
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		response.Error(w, http.StatusInternalServerError, "cannot write file")
		return
	}

	userIDs := findUserIDs(r.Context(), api, meta.Grant)
	if err := api.Service.AddDocumentAccess(r.Context(), docID, userIDs); err != nil {
		response.Error(w, http.StatusInternalServerError, "cannot save document access")
		return
	}

	jsonData := parseOptionalJSON(r, api)

	response.JSON(w, http.StatusOK, response.DocumentUploadResponse{
		Data: response.DocumentUploadData{
			File: originalName,
			JSON: jsonData,
		},
	})
}

func (api *DocumentHandler) getDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Логика получения списка документов
}

func (api *DocumentHandler) getDocumentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "missing document id")
		return
	}

	doc, err := api.Service.GetDocumentByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "document not found")
		return
	}

	token := r.URL.Query().Get("token")

	fmt.Println(token)

	if token == "" {
		response.Error(w, http.StatusUnauthorized, "missing token")
		return
	}
	userID, err := api.Service.ValidateToken(r.Context(), token)
	if err != nil {
		response.Error(w, http.StatusForbidden, "invalid token")
		return
	}

	if !canUserAccessDocument(r.Context(), api, doc, userID) {
		response.Error(w, http.StatusForbidden, "access denied")
		return
	}

	if doc.HasFile {
		filePath := filepath.Join("data/uploads", doc.FileName)
		cleanPath := filepath.Clean(filePath)

		if rel, err := filepath.Rel("data/uploads", cleanPath); err != nil || strings.HasPrefix(rel, "..") {
			response.Error(w, http.StatusForbidden, "invalid file path")
			return
		}

		w.Header().Set("Content-Type", doc.MimeType)
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, doc.FileName))
		http.ServeFile(w, r, cleanPath)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": doc,
	})
}

func (api *DocumentHandler) deleteDocumentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "missing document id")
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		response.Error(w, http.StatusUnauthorized, "missing token")
		return
	}
	userID, err := api.Service.ValidateToken(r.Context(), token)
	if err != nil {
		response.Error(w, http.StatusForbidden, "invalid token")
		return
	}

	doc, err := api.Service.GetDocumentByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "document not found")
		return
	}

	if doc.OwnerID != userID {
		response.Error(w, http.StatusForbidden, "access denied")
		return
	}

	err = api.Service.DeleteDocument(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "cannot delete document")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"response": map[string]bool{
			id: true,
		},
	})
}

func findUserIDs(ctx context.Context, api *DocumentHandler, logins []string) []string {
	var userIDs []string
	for _, login := range logins {
		if userID, err := api.Service.FindUserIDByLogin(ctx, login); err == nil {
			userIDs = append(userIDs, userID)
		}
	}
	return userIDs
}

func parseOptionalJSON(r *http.Request, api *DocumentHandler) map[string]interface{} {
	jsonStr := r.FormValue("json")
	if jsonStr == "" {
		return nil
	}
	data, _ := api.Service.ParseOptionalJSON(jsonStr)
	return data
}

func canUserAccessDocument(ctx context.Context, api *DocumentHandler, doc *domain.Document, userID string) bool {
	if doc.IsPublic {
		return true
	}
	if userID != "" && doc.OwnerID == userID {
		return true
	}
	if userID != "" {
		ok, _ := api.Service.HasDocumentAccess(ctx, doc.ID, userID)
		return ok
	}
	return false
}

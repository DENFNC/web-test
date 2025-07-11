package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

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

	fileHeader, err := utils.SaveUploadedFile(r)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	docID, err := api.Service.SaveDocument(r.Context(), *meta, ownerID)
	if err != nil {
		fmt.Println(err)
		response.Error(w, http.StatusInternalServerError, "cannot save document")
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
			File: fileHeader.Filename,
			JSON: jsonData,
		},
	})
}

func (api *DocumentHandler) getDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Логика получения списка документов
}

func (api *DocumentHandler) getDocumentHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Логика получения конкретного документа
}

func (api *DocumentHandler) deleteDocumentHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Логика удаления документа
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

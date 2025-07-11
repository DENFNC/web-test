package service

import (
	"context"
	"encoding/json"

	"github.com/DENFNC/web-test/internal/domain"
	"github.com/DENFNC/web-test/internal/infra/psql/repository"
	"github.com/DENFNC/web-test/internal/transport/dto/request"
	"github.com/google/uuid"
)

type DocumentService struct {
	DocRepo  *repository.DocumentRepository
	AuthRepo *repository.AuthRepository
}

func NewDocumentService(docRepo *repository.DocumentRepository, authRepo *repository.AuthRepository) *DocumentService {
	return &DocumentService{
		DocRepo:  docRepo,
		AuthRepo: authRepo,
	}
}

func (s *DocumentService) ValidateToken(ctx context.Context, token string) (string, error) {
	return s.AuthRepo.GetUserIDByToken(ctx, token)
}

func (s *DocumentService) FindUserIDByLogin(ctx context.Context, login string) (string, error) {
	return s.AuthRepo.GetUserIDByLogin(ctx, login)
}

func (s *DocumentService) SaveDocument(ctx context.Context, meta request.DocumentMetaRequest, ownerID string) (string, error) {
	id := uuid.New()

	doc := &domain.Document{
		ID:       id.String(),
		FileName: meta.Name,
		MimeType: meta.Mime,
		HasFile:  true,
		IsPublic: meta.Public,
		OwnerID:  ownerID,
	}
	return s.DocRepo.SaveDocument(ctx, doc)
}

func (s *DocumentService) AddDocumentAccess(ctx context.Context, documentID string, userIDs []string) error {
	for _, userID := range userIDs {
		if err := s.DocRepo.AddDocumentAccess(ctx, documentID, userID); err != nil {
			return err
		}
	}
	return nil
}

func (s *DocumentService) ParseOptionalJSON(jsonStr string) (map[string]interface{}, error) {
	if jsonStr == "" {
		return nil, nil
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, err
	}
	return data, nil
}

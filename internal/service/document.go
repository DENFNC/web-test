package service

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/DENFNC/web-test/internal/model"
	"github.com/DENFNC/web-test/internal/repository"
)

type DocumentService struct {
	repo repository.Repo
}

func NewDocumentService(repo repository.Repo) *DocumentService {
	return &DocumentService{repo: repo}
}

// Загрузка документа
func (s *DocumentService) Upload(doc *model.Document) error {
	if doc.ID == "" {
		doc.ID = genDocID()
	}
	doc.Created = time.Now()
	s.repo.AddDocument(doc)
	return nil
}

// Удаление документа
func (s *DocumentService) Delete(id, login string) error {
	doc, ok := s.repo.GetDocument(id)
	if !ok {
		return errors.New("not found")
	}
	if doc.Owner != login {
		return errors.New("forbidden")
	}
	s.repo.DeleteDocument(id)
	return nil
}

// Получение одного документа (с учётом прав)
func (s *DocumentService) Get(id, login string) (*model.Document, error) {
	doc, ok := s.repo.GetDocument(id)
	if !ok {
		return nil, errors.New("not found")
	}
	if !doc.Public && doc.Owner != login && !contains(doc.Grant, login) {
		return nil, errors.New("forbidden")
	}
	return doc, nil
}

// Получение списка документов (фильтрация, сортировка, права)
func (s *DocumentService) List(login, filterKey, filterValue string, limit int) []*model.Document {
	docs := s.repo.ListDocuments()
	var result []*model.Document
	for _, d := range docs {
		if !d.Public && d.Owner != login && !contains(d.Grant, login) {
			continue
		}
		if filterKey != "" && filterValue != "" {
			v := getFieldValue(d, filterKey)
			if v != filterValue {
				continue
			}
		}
		result = append(result, d)
	}
	// Сортировка по имени и дате
	sort.Slice(result, func(i, j int) bool {
		if result[i].Name == result[j].Name {
			return result[i].Created.Before(result[j].Created)
		}
		return result[i].Name < result[j].Name
	})
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result
}

// --- Вспомогательные функции ---
func contains(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func getFieldValue(d *model.Document, key string) string {
	switch key {
	case "id":
		return d.ID
	case "name":
		return d.Name
	case "mime":
		return d.Mime
	case "owner":
		return d.Owner
	}
	return ""
}

func genDocID() string {
	return strings.ReplaceAll(time.Now().Format("20060102150405.000000"), ".", "")
}

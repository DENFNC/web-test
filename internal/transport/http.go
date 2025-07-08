package transport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/DENFNC/web-test/internal/model"
	"github.com/DENFNC/web-test/internal/repository"
	"github.com/DENFNC/web-test/internal/service"
)

type cacheEntry struct {
	data   []byte
	header http.Header
	status int
	expiry time.Time
}

type Server struct {
	docSvc  *service.DocumentService
	userSvc *service.UserService
	repo    repository.Repo
	cache   repository.Cache
	cacheMu sync.RWMutex
}

func NewServerWithCache(repo repository.Repo, docSvc *service.DocumentService, userSvc *service.UserService, cache repository.Cache) *Server {
	return &Server{
		docSvc:  docSvc,
		userSvc: userSvc,
		repo:    repo,
		cache:   cache,
	}
}

// --- Кэширование ---
func (s *Server) getCache(key string) ([]byte, bool) {
	return s.cache.Get(key)
}

func (s *Server) setCache(key string, data []byte) {
	s.cache.Set(key, data, 2*time.Minute)
}

func (s *Server) invalidateCache(keys ...string) {
	s.cache.Invalidate(keys...)
}

// --- Форматированный ответ ---
func writeResponse(w http.ResponseWriter, status int, errObj any, respObj any, dataObj any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	m := make(map[string]any)
	if errObj != nil {
		m["error"] = errObj
	}
	if respObj != nil {
		m["response"] = respObj
	}
	if dataObj != nil {
		m["data"] = dataObj
	}
	json.NewEncoder(w).Encode(m)
}

// --- Регистрация ---
func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeResponse(w, 405, map[string]any{"code": 405, "text": "method not allowed"}, nil, nil)
		return
	}
	var req struct {
		Token string `json:"token"`
		Login string `json:"login"`
		Pswd  string `json:"pswd"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeResponse(w, 400, map[string]any{"code": 400, "text": "bad request"}, nil, nil)
		return
	}
	err := s.userSvc.Register(req.Token, req.Login, req.Pswd)
	if err != nil {
		code := 400
		if err.Error() == "forbidden" {
			code = 403
		}
		writeResponse(w, code, map[string]any{"code": code, "text": err.Error()}, nil, nil)
		return
	}
	writeResponse(w, 200, nil, map[string]any{"login": req.Login}, nil)
}

// --- Аутентификация ---
func (s *Server) handleAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeResponse(w, 405, map[string]any{"code": 405, "text": "method not allowed"}, nil, nil)
		return
	}
	if err := r.ParseForm(); err != nil {
		writeResponse(w, 400, map[string]any{"code": 400, "text": "bad request"}, nil, nil)
		return
	}
	login := r.FormValue("login")
	pswd := r.FormValue("pswd")
	token, err := s.userSvc.Auth(login, pswd)
	if err != nil {
		writeResponse(w, 401, map[string]any{"code": 401, "text": "unauthorized"}, nil, nil)
		return
	}
	writeResponse(w, 200, nil, map[string]any{"token": token}, nil)
}

// --- Завершение сессии ---
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeResponse(w, 405, map[string]any{"code": 405, "text": "method not allowed"}, nil, nil)
		return
	}
	token := strings.TrimPrefix(r.URL.Path, "/api/auth/")
	s.userSvc.Logout(token)
	writeResponse(w, 200, nil, map[string]any{token: true}, nil)
	s.invalidateCache() // Инвалидация всего кеша для простоты
}

// --- Загрузка документа ---
func (s *Server) handleUploadDoc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeResponse(w, 405, map[string]any{"code": 405, "text": "method not allowed"}, nil, nil)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeResponse(w, 400, map[string]any{"code": 400, "text": "bad multipart"}, nil, nil)
		return
	}
	metaStr := r.FormValue("meta")
	var meta struct {
		Name   string   `json:"name"`
		File   bool     `json:"file"`
		Public bool     `json:"public"`
		Token  string   `json:"token"`
		Mime   string   `json:"mime"`
		Grant  []string `json:"grant"`
	}
	if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
		writeResponse(w, 400, map[string]any{"code": 400, "text": "bad meta"}, nil, nil)
		return
	}
	login, ok := s.userSvc.CheckToken(meta.Token)
	if !ok {
		writeResponse(w, 401, map[string]any{"code": 401, "text": "unauthorized"}, nil, nil)
		return
	}
	var jsonData any
	if jsonStr := r.FormValue("json"); jsonStr != "" {
		json.Unmarshal([]byte(jsonStr), &jsonData)
	}
	var fileName string
	if meta.File {
		file, header, err := r.FormFile("file")
		if err != nil {
			writeResponse(w, 400, map[string]any{"code": 400, "text": "no file"}, nil, nil)
			return
		}
		defer file.Close()
		fileName = header.Filename
		// Файл не сохраняем, только метаданные (для теста)
	}
	doc := &model.Document{
		Name:   meta.Name,
		File:   meta.File,
		Public: meta.Public,
		Mime:   meta.Mime,
		Grant:  meta.Grant,
		Owner:  login,
		JSON:   jsonData,
	}
	err := s.docSvc.Upload(doc)
	if err != nil {
		writeResponse(w, 500, map[string]any{"code": 500, "text": "internal error"}, nil, nil)
		return
	}
	writeResponse(w, 200, nil, nil, map[string]any{"json": jsonData, "file": fileName})
	s.invalidateCache() // Инвалидация всего кеша для простоты
}

// --- Удаление документа ---
func (s *Server) handleDeleteDoc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeResponse(w, 405, map[string]any{"code": 405, "text": "method not allowed"}, nil, nil)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/docs/")
	token := r.URL.Query().Get("token")
	login, ok := s.userSvc.CheckToken(token)
	if !ok {
		writeResponse(w, 401, map[string]any{"code": 401, "text": "unauthorized"}, nil, nil)
		return
	}
	err := s.docSvc.Delete(id, login)
	if err != nil {
		code := 403
		if err.Error() == "not found" {
			code = 404
		}
		writeResponse(w, code, map[string]any{"code": code, "text": err.Error()}, nil, nil)
		return
	}
	writeResponse(w, 200, nil, map[string]any{id: true}, nil)
	s.invalidateCache() // Инвалидация всего кеша для простоты
}

// --- Получение списка документов ---
func (s *Server) handleListDocs(w http.ResponseWriter, r *http.Request) {
	cacheKey := r.Method + r.URL.RequestURI()
	// --- Получение из кэша ---
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		if buf, ok := s.getCache(cacheKey); ok {
			w.WriteHeader(200)
			w.Write(buf)
			return
		}
	}
	token := r.URL.Query().Get("token")
	login, _ := s.userSvc.CheckToken(token)
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	docs := s.docSvc.List(login, key, value, limit)
	resp := map[string]any{"docs": docs}
	buf, _ := json.Marshal(map[string]any{"data": resp})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(buf)
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		s.setCache(cacheKey, buf)
	}
}

// --- Получение одного документа ---
func (s *Server) handleGetDoc(w http.ResponseWriter, r *http.Request) {
	cacheKey := r.Method + r.URL.RequestURI()
	// --- Получение из кэша ---
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		if buf, ok := s.getCache(cacheKey); ok {
			w.WriteHeader(200)
			w.Write(buf)
			return
		}
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/docs/")
	token := r.URL.Query().Get("token")
	login, _ := s.userSvc.CheckToken(token)
	doc, err := s.docSvc.Get(id, login)
	if err != nil {
		code := 403
		if err.Error() == "not found" {
			code = 404
		}
		writeResponse(w, code, map[string]any{"code": code, "text": err.Error()}, nil, nil)
		return
	}
	if doc.File {
		w.Header().Set("Content-Type", doc.Mime)
		w.Header().Set("Content-Disposition", "attachment; filename="+doc.Name)
		w.WriteHeader(200)
		w.Write([]byte("FAKE FILE DATA")) // Для теста: отдаём заглушку
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			buf := []byte("FAKE FILE DATA")
			s.setCache(cacheKey, buf)
		}
		return
	}
	buf, _ := json.Marshal(map[string]any{"data": doc.JSON})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(buf)
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		s.setCache(cacheKey, buf)
	}
}

// --- Роутер ---
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/api/register":
		s.handleRegister(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/api/auth":
		s.handleAuth(w, r)
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/api/auth/"):
		s.handleLogout(w, r)
	case (r.Method == http.MethodGet || r.Method == http.MethodHead) && r.URL.Path == "/api/docs":
		s.handleListDocs(w, r)
	case (r.Method == http.MethodGet || r.Method == http.MethodHead) && strings.HasPrefix(r.URL.Path, "/api/docs/"):
		s.handleGetDoc(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/api/docs":
		s.handleUploadDoc(w, r)
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/api/docs/"):
		s.handleDeleteDoc(w, r)
	default:
		writeResponse(w, 405, map[string]any{"code": 405, "text": "method not allowed"}, nil, nil)
	}
}

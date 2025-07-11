package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Body json.RawMessage `json:"response"`
}

type ErrorResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload == nil {
		return
	}

	data, err := toJSON(payload)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := Response{
		Body: data,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func Error(w http.ResponseWriter, status int, message string) {
	resp := ErrorResponse{
		Code: status,
		Text: message,
	}
	JSON(w, status, resp)
}

func toJSON(v any) (json.RawMessage, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

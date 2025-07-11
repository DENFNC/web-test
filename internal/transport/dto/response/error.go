package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Body json.RawMessage `json:"response"`
}

type ErrorResponse struct {
	Error struct {
		Code int    `json:"code"`
		Text string `json:"text"`
	} `json:"error"`
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{}
	resp.Error.Code = status
	resp.Error.Text = message

	_ = json.NewEncoder(w).Encode(resp)
}

func toJSON(v any) (json.RawMessage, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

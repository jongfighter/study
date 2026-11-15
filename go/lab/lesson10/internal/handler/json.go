package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// writeJSON은 JSON 응답을 씁니다. 헤더를 이미 보낸 뒤에는 인코딩 에러를 할 수 있는 일이 없어 무시합니다.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func decodeJSON(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func parseID(r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		return 0, errors.New("id는 1 이상의 정수여야 합니다")
	}
	return id, nil
}

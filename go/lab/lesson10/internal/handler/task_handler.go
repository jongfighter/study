// Package handler는 HTTP 요청을 읽고 서비스를 호출해 HTTP 응답으로 바꿉니다.
package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"golab/lesson10/internal/domain"
)

// TaskService는 핸들러가 필요로 하는 서비스 동작입니다.
// 인터페이스를 사용하는 쪽(handler)에서 정의하므로 테스트에서 가짜 구현을 끼우기 쉽습니다.
type TaskService interface {
	List(ctx context.Context) ([]domain.Task, error)
	Create(ctx context.Context, title string) (domain.Task, error)
	Get(ctx context.Context, id int) (domain.Task, error)
	Update(ctx context.Context, id int, in domain.UpdateInput) (domain.Task, error)
	Delete(ctx context.Context, id int) error
}

// Handler는 Task API의 HTTP 핸들러 모음입니다.
type Handler struct {
	tasks  TaskService
	logger *slog.Logger
}

// New는 tasks와 logger를 주입받아 Handler를 만듭니다.
func New(tasks TaskService, logger *slog.Logger) *Handler {
	return &Handler{tasks: tasks, logger: logger}
}

// Routes는 경로와 핸들러를 연결한 라우터를 돌려줍니다.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /tasks", h.listTasks)
	mux.HandleFunc("POST /tasks", h.createTask)
	mux.HandleFunc("GET /tasks/{id}", h.getTask)
	mux.HandleFunc("PATCH /tasks/{id}", h.updateTask)
	mux.HandleFunc("DELETE /tasks/{id}", h.deleteTask)
	return mux
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.tasks.List(r.Context())
	if err != nil {
		h.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

type createTaskRequest struct {
	Title string `json:"title"`
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "요청 본문을 읽을 수 없습니다: "+err.Error())
		return
	}

	task, err := h.tasks.Create(r.Context(), req.Title)
	if err != nil {
		h.writeServiceError(w, r, err)
		return
	}
	w.Header().Set("Location", "/tasks/"+strconv.Itoa(task.ID))
	writeJSON(w, http.StatusCreated, task)
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.tasks.Get(r.Context(), id)
	if err != nil {
		h.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

type updateTaskRequest struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var req updateTaskRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "요청 본문을 읽을 수 없습니다: "+err.Error())
		return
	}

	task, err := h.tasks.Update(r.Context(), id, domain.UpdateInput{Title: req.Title, Done: req.Done})
	if err != nil {
		h.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.tasks.Delete(r.Context(), id); err != nil {
		h.writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// writeServiceError는 서비스가 돌려준 에러를 HTTP 상태 코드로 바꿉니다.
func (h *Handler) writeServiceError(w http.ResponseWriter, r *http.Request, err error) {
	var validation *domain.ValidationError
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.As(err, &validation):
		writeError(w, http.StatusBadRequest, validation.Message)
	default:
		h.logger.ErrorContext(r.Context(), "처리하지 못한 에러", "error", err, "request_id", RequestIDFrom(r.Context()))
		writeError(w, http.StatusInternalServerError, "서버 내부 오류")
	}
}

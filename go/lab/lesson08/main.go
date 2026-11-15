package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Task는 할 일 하나를 나타냅니다.
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// taskStore는 Task를 메모리에 보관합니다.
// net/http는 요청마다 고루틴을 따로 실행하므로 Mutex로 보호합니다.
type taskStore struct {
	mu     sync.Mutex
	nextID int
	tasks  map[int]Task
}

func newTaskStore() *taskStore {
	return &taskStore{nextID: 1, tasks: make(map[int]Task)}
}

func (s *taskStore) list() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (s *taskStore) create(title string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := Task{ID: s.nextID, Title: title, CreatedAt: time.Now().UTC()}
	s.nextID++
	s.tasks[t.ID] = t
	return t
}

func (s *taskStore) get(id int) (Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	return t, ok
}

func (s *taskStore) update(id int, title *string, done *bool) (Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	if !ok {
		return Task{}, false
	}
	if title != nil {
		t.Title = *title
	}
	if done != nil {
		t.Done = *done
	}
	s.tasks[id] = t
	return t, true
}

func (s *taskStore) remove(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.tasks[id]
	delete(s.tasks, id)
	return ok
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("응답 쓰기 실패: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func decodeJSON(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

type server struct {
	store *taskStore
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /tasks", s.listTasks)
	mux.HandleFunc("POST /tasks", s.createTask)
	mux.HandleFunc("GET /tasks/{id}", s.getTask)
	mux.HandleFunc("PATCH /tasks/{id}", s.updateTask)
	mux.HandleFunc("DELETE /tasks/{id}", s.deleteTask)
	return mux
}

func (s *server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) listTasks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.list())
}

type createTaskRequest struct {
	Title string `json:"title"`
}

func (s *server) createTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "요청 본문을 읽을 수 없습니다: "+err.Error())
		return
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		writeError(w, http.StatusBadRequest, "title은 비어 있을 수 없습니다")
		return
	}

	task := s.store.create(title)
	w.Header().Set("Location", "/tasks/"+strconv.Itoa(task.ID))
	writeJSON(w, http.StatusCreated, task)
}

func parseID(r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		return 0, errors.New("id는 1 이상의 정수여야 합니다")
	}
	return id, nil
}

func (s *server) getTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	task, ok := s.store.get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "할 일을 찾을 수 없습니다")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

type updateTaskRequest struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

func (s *server) updateTask(w http.ResponseWriter, r *http.Request) {
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
	if req.Title == nil && req.Done == nil {
		writeError(w, http.StatusBadRequest, "수정할 필드(title, done)가 없습니다")
		return
	}
	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			writeError(w, http.StatusBadRequest, "title은 비어 있을 수 없습니다")
			return
		}
		req.Title = &title
	}

	task, ok := s.store.update(id, req.Title, req.Done)
	if !ok {
		writeError(w, http.StatusNotFound, "할 일을 찾을 수 없습니다")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (s *server) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !s.store.remove(id) {
		writeError(w, http.StatusNotFound, "할 일을 찾을 수 없습니다")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &server{store: newTaskStore()}
	log.Printf("Task API 서버 시작: http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, srv.routes()))
}

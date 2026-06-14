package handler

import (
	"io"
	"net/http"

	"github.com/noedaka/go-config-parser/internal/parser"
	"github.com/noedaka/go-config-parser/internal/service"
)

type Handler struct {
	rules  []service.Rule
	parser parser.Parser
}

func NewHandler(rules []service.Rule, parser parser.Parser) *Handler {
	return &Handler{rules: rules, parser: parser}
}

func (h *Handler) ConfigRecommendationsByFileHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("config")
	if err != nil {
		http.Error(w, "Error getting the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "File reading error", http.StatusInternalServerError)
		return
	}

	data, err := h.parser.ParseConfig(fileBytes)
	if err != nil {
		http.Error(w, "Cannot parse data", http.StatusInternalServerError)
		return
	}

	var issues []service.Issue
	for _, rule := range h.rules {
		issues = append(issues, rule.Check(data)...)
	}

	response := service.FormatIssues(issues)

	if len(issues) > 0 {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(response))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Рекомендации к конфигурации не требуются"))
}

func (h *Handler) ConfigRecommendationsByBodyHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	data, err := h.parser.ParseConfig(body)
	if err != nil {
		http.Error(w, "Cannot parse data", http.StatusInternalServerError)
		return
	}

	var issues []service.Issue
	for _, rule := range h.rules {
		issues = append(issues, rule.Check(data)...)
	}

	response := service.FormatIssues(issues)

	if len(issues) > 0 {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(response))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Рекомендации к конфигурации не требуются"))
}

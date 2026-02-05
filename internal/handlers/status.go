package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yatoenough/moniserve/internal/checker"
)

type StatusHandler struct {
	checker *checker.Checker
}

func NewStatusHandler(ch *checker.Checker) *StatusHandler {
	return &StatusHandler{
		checker: ch,
	}
}

func (sh *StatusHandler) Handle(w http.ResponseWriter, r *http.Request) {
	results := sh.checker.CheckAll(r.Context())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

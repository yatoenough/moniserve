package handlers

import (
	"net/http"
)

func Setup(m *http.ServeMux) {
	m.HandleFunc("/", hello)
}

package handler

import (
	"net/http"
	"time"
)

type HomeHandler struct{}

func (h *HomeHandler) Index(c *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templateName := p + "web/template/index.html"
		data := &TemplateData{
			Now: time.Now(),
		}
		c.render(w, r, nil, templateName, data)
	}
}

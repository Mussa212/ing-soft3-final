package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddlewareSetsHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodOptions, "/any", nil)

	CORSMiddleware()(c)

	headers := w.Result().Header
	if headers.Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("missing allow origin header")
	}
	if headers.Get("Access-Control-Allow-Headers") == "" || headers.Get("Access-Control-Allow-Methods") == "" {
		t.Fatalf("missing allow headers/methods")
	}
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 for preflight, got %d", w.Code)
	}
}

package request

import (
	"io"
	"testing"
)

func TestFactory_Build(t *testing.T) {
	factory := NewFactory("GET", "http://localhost:8080/users/{{.id}}?name={{.name}}", "", "")
	header := []string{"id", "name"}
	row := []string{"1", "gopher"}
	req, err := factory.Build(header, row)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != "GET" {
		t.Errorf("expected method GET, but got %s", req.Method)
	}
	expectedURL := "http://localhost:8080/users/1?name=gopher"
	if req.URL.String() != expectedURL {
		t.Errorf("expected url %s, but got %s", expectedURL, req.URL.String())
	}
}

func TestFactory_Build_WithHeaderAndBody(t *testing.T) {
	factory := NewFactory(
		"POST",
		"http://localhost:8080/users",
		"Content-Type: application/json\nX-Request-ID: {{.request_id}}",
		`{"name": "{{.name}}"}`,
	)
	header := []string{"request_id", "name"}
	row := []string{"xyz-123", "gopher"}
	req, err := factory.Build(header, row)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Method != "POST" {
		t.Errorf("expected method POST, but got %s", req.Method)
	}

	// Check header
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type 'application/json', but got '%s'", req.Header.Get("Content-Type"))
	}
	if req.Header.Get("X-Request-ID") != "xyz-123" {
		t.Errorf("expected X-Request-ID 'xyz-123', but got '%s'", req.Header.Get("X-Request-ID"))
	}

	// Check body
	expectedBody := `{"name": "gopher"}`
	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != expectedBody {
		t.Errorf("expected body '%s', but got '%s'", expectedBody, string(body))
	}
}
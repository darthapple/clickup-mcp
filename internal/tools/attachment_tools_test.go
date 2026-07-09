package tools

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupCreateTaskAttachment(t *testing.T) {
	t.Run("required arg validation", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterAttachmentTools(s, c)
		res := callTool(t, s, "clickup_create_task_attachment", map[string]any{"task_id": "t1"})
		if !res.IsError {
			t.Error("IsError = false, want true (missing file_path)")
		}
		if hit {
			t.Error("handler hit the fake server despite missing required arg")
		}
	})

	t.Run("argument wiring uploads multipart file", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, "hello.txt")
		if err := os.WriteFile(filePath, []byte("hello world"), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}

		var gotMethod, gotPath, gotContentType string
		var gotFileName string
		var gotFileContents []byte
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotContentType = r.Header.Get("Content-Type")
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				t.Errorf("ParseMultipartForm: %v", err)
			} else {
				f, hdr, err := r.FormFile("attachment")
				if err != nil {
					t.Errorf("FormFile: %v", err)
				} else {
					defer f.Close()
					gotFileName = hdr.Filename
					buf := make([]byte, 1024)
					n, _ := f.Read(buf)
					gotFileContents = buf[:n]
				}
			}
			_, _ = w.Write([]byte(`{"id":"att1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterAttachmentTools(s, c)
		res := callTool(t, s, "clickup_create_task_attachment", map[string]any{
			"task_id":   "t1",
			"file_path": filePath,
		})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/task/t1/attachment" {
			t.Errorf("path = %q, want /task/t1/attachment", gotPath)
		}
		if !strings.HasPrefix(gotContentType, "multipart/form-data; boundary=") {
			t.Errorf("Content-Type = %q, want multipart/form-data; boundary=...", gotContentType)
		}
		if gotFileName != "hello.txt" {
			t.Errorf("uploaded filename = %q, want hello.txt", gotFileName)
		}
		if string(gotFileContents) != "hello world" {
			t.Errorf("uploaded contents = %q, want %q", gotFileContents, "hello world")
		}
	})

	t.Run("missing local file surfaces as error", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterAttachmentTools(s, c)
		res := callTool(t, s, "clickup_create_task_attachment", map[string]any{
			"task_id":   "t1",
			"file_path": "/nonexistent/path/does-not-exist.txt",
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		if hit {
			t.Error("handler hit the fake server despite unreadable file")
		}
	})

	t.Run("error passthrough", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, "hello.txt")
		if err := os.WriteFile(filePath, []byte("hello world"), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterAttachmentTools(s, c)
		res := callTool(t, s, "clickup_create_task_attachment", map[string]any{
			"task_id":   "t1",
			"file_path": filePath,
		})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("text = %q, want %q", textOf(t, res), want)
		}
	})
}

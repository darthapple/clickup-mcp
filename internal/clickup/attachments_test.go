package clickup

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateTaskAttachmentUploadsMultipartFile(t *testing.T) {
	var gotMethod, gotPath, gotContentType, gotFilename, gotContent string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotContentType = r.Header.Get("Content-Type")

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Errorf("ParseMultipartForm: %v", err)
		}
		file, header, err := r.FormFile("attachment")
		if err != nil {
			t.Fatalf("FormFile: %v", err)
		}
		defer file.Close()
		gotFilename = header.Filename
		b, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("read file: %v", err)
		}
		gotContent = string(b)

		_, _ = w.Write([]byte(`{"id":"att1","title":"hello.txt"}`))
	})

	dir := t.TempDir()
	filePath := filepath.Join(dir, "hello.txt")
	if err := os.WriteFile(filePath, []byte("hello world"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	out, err := c.CreateTaskAttachment(context.Background(), "task1", filePath)
	if err != nil {
		t.Fatalf("CreateTaskAttachment: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/task/task1/attachment" {
		t.Errorf("path = %q, want /task/task1/attachment", gotPath)
	}
	if !strings.HasPrefix(gotContentType, "multipart/form-data; boundary=") {
		t.Errorf("Content-Type = %q, want multipart/form-data; boundary=...", gotContentType)
	}
	if gotFilename != "hello.txt" {
		t.Errorf("uploaded filename = %q, want hello.txt", gotFilename)
	}
	if gotContent != "hello world" {
		t.Errorf("uploaded content = %q, want %q", gotContent, "hello world")
	}
	m, ok := out.(map[string]any)
	if !ok || m["id"] != "att1" {
		t.Errorf("decoded = %+v", out)
	}
}

func TestCreateTaskAttachmentMissingFileReturnsError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("server should not be called when the local file can't be opened")
	})

	_, err := c.CreateTaskAttachment(context.Background(), "task1", filepath.Join(t.TempDir(), "does-not-exist.txt"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAttachmentsAPIErrorReturnsAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		_, _ = w.Write([]byte(`{"err":"File too large","ECODE":"ATTACH_001"}`))
	})

	dir := t.TempDir()
	filePath := filepath.Join(dir, "hello.txt")
	if err := os.WriteFile(filePath, []byte("hello world"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := c.CreateTaskAttachment(context.Background(), "task1", filePath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusRequestEntityTooLarge || apiErr.ECode != "ATTACH_001" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}

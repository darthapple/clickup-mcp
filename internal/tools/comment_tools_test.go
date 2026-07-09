package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestClickupListComments(t *testing.T) {
	cases := []struct {
		name       string
		tool       string
		idArg      string
		idVal      string
		wantPath   string
		missingArg string
	}{
		{"task", "clickup_list_task_comments", "task_id", "task1", "/task/task1/comment", "task_id"},
		{"list", "clickup_list_list_comments", "list_id", "list1", "/list/list1/comment", "list_id"},
		{"view", "clickup_list_view_comments", "view_id", "view1", "/view/view1/comment", "view_id"},
	}

	for _, tc := range cases {
		t.Run(tc.name+"/missing_required_arg", func(t *testing.T) {
			hit := false
			c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				hit = true
				t.Errorf("server should not have been called")
			})
			s := server.NewMCPServer("test", "1.0.0")
			RegisterCommentTools(s, c)
			res := callTool(t, s, tc.tool, map[string]any{})
			if !res.IsError {
				t.Error("IsError = false, want true")
			}
			if hit {
				t.Error("server was hit despite missing required arg")
			}
		})

		t.Run(tc.name+"/wiring", func(t *testing.T) {
			c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("method = %s, want GET", r.Method)
				}
				if r.URL.Path != tc.wantPath {
					t.Errorf("path = %s, want %s", r.URL.Path, tc.wantPath)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"comments":[{"id":"c1"}]}`))
			})
			s := server.NewMCPServer("test", "1.0.0")
			RegisterCommentTools(s, c)
			res := callTool(t, s, tc.tool, map[string]any{tc.idArg: tc.idVal})
			if res.IsError {
				t.Fatalf("IsError = true, want false: %s", textOf(t, res))
			}
			if !strings.Contains(textOf(t, res), "c1") {
				t.Errorf("result = %q, want it to contain c1", textOf(t, res))
			}
		})

		t.Run(tc.name+"/error_passthrough", func(t *testing.T) {
			c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
			})
			s := server.NewMCPServer("test", "1.0.0")
			RegisterCommentTools(s, c)
			res := callTool(t, s, tc.tool, map[string]any{tc.idArg: tc.idVal})
			if !res.IsError {
				t.Fatal("IsError = false, want true")
			}
			want := "ClickUp API error 404 [X_001]: not found"
			if textOf(t, res) != want {
				t.Errorf("error text = %q, want %q", textOf(t, res), want)
			}
		})
	}
}

func TestClickupCreateComment(t *testing.T) {
	cases := []struct {
		name     string
		tool     string
		idArg    string
		idVal    string
		wantPath string
	}{
		{"task", "clickup_create_task_comment", "task_id", "task1", "/task/task1/comment"},
		{"list", "clickup_create_list_comment", "list_id", "list1", "/list/list1/comment"},
		{"view", "clickup_create_view_comment", "view_id", "view1", "/view/view1/comment"},
		{"reply", "clickup_create_comment_reply", "comment_id", "comment1", "/comment/comment1/reply"},
	}

	for _, tc := range cases {
		t.Run(tc.name+"/missing_comment_text", func(t *testing.T) {
			hit := false
			c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				hit = true
				t.Errorf("server should not have been called")
			})
			s := server.NewMCPServer("test", "1.0.0")
			RegisterCommentTools(s, c)
			res := callTool(t, s, tc.tool, map[string]any{tc.idArg: tc.idVal})
			if !res.IsError {
				t.Error("IsError = false, want true")
			}
			if hit {
				t.Error("server was hit despite missing required comment_text")
			}
		})

		t.Run(tc.name+"/wiring", func(t *testing.T) {
			c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("method = %s, want POST", r.Method)
				}
				if r.URL.Path != tc.wantPath {
					t.Errorf("path = %s, want %s", r.URL.Path, tc.wantPath)
				}
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("decode body: %v", err)
				}
				if body["comment_text"] != "hello" {
					t.Errorf("comment_text = %v, want hello", body["comment_text"])
				}
				if body["notify_all"] != true {
					t.Errorf("notify_all = %v, want true", body["notify_all"])
				}
				if body["assignee"] != "user1" {
					t.Errorf("assignee = %v, want user1", body["assignee"])
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"id":"comment42"}`))
			})
			s := server.NewMCPServer("test", "1.0.0")
			RegisterCommentTools(s, c)
			res := callTool(t, s, tc.tool, map[string]any{
				tc.idArg:       tc.idVal,
				"comment_text": "hello",
				"notify_all":   true,
				"assignee":     "user1",
			})
			if res.IsError {
				t.Fatalf("IsError = true, want false: %s", textOf(t, res))
			}
			if !strings.Contains(textOf(t, res), "comment42") {
				t.Errorf("result = %q, want it to contain comment42", textOf(t, res))
			}
		})

		t.Run(tc.name+"/only_required_field_in_body", func(t *testing.T) {
			c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("decode body: %v", err)
				}
				if len(body) != 1 {
					t.Errorf("body = %v, want only comment_text set", body)
				}
				if body["comment_text"] != "hi" {
					t.Errorf("comment_text = %v, want hi", body["comment_text"])
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"id":"comment1"}`))
			})
			s := server.NewMCPServer("test", "1.0.0")
			RegisterCommentTools(s, c)
			res := callTool(t, s, tc.tool, map[string]any{tc.idArg: tc.idVal, "comment_text": "hi"})
			if res.IsError {
				t.Fatalf("IsError = true, want false: %s", textOf(t, res))
			}
		})
	}
}

func TestClickupUpdateComment(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCommentTools(s, c)
		res := callTool(t, s, "clickup_update_comment", map[string]any{"comment_text": "x"})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing comment_id")
		}
	})

	t.Run("partial_update_only_resolved", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("method = %s, want PUT", r.Method)
			}
			if r.URL.Path != "/comment/c1" {
				t.Errorf("path = %s, want /comment/c1", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if len(body) != 1 {
				t.Errorf("body = %v, want only resolved set", body)
			}
			if body["resolved"] != true {
				t.Errorf("resolved = %v, want true", body["resolved"])
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"c1"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCommentTools(s, c)
		res := callTool(t, s, "clickup_update_comment", map[string]any{"comment_id": "c1", "resolved": true})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
	})

	t.Run("error_passthrough", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"not found","ECODE":"X_001"}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCommentTools(s, c)
		res := callTool(t, s, "clickup_update_comment", map[string]any{"comment_id": "c1", "comment_text": "x"})
		if !res.IsError {
			t.Fatal("IsError = false, want true")
		}
		want := "ClickUp API error 404 [X_001]: not found"
		if textOf(t, res) != want {
			t.Errorf("error text = %q, want %q", textOf(t, res), want)
		}
	})
}

func TestClickupDeleteComment(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCommentTools(s, c)
		res := callTool(t, s, "clickup_delete_comment", map[string]any{})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing comment_id")
		}
	})

	t.Run("wiring", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("method = %s, want DELETE", r.Method)
			}
			if r.URL.Path != "/comment/c1" {
				t.Errorf("path = %s, want /comment/c1", r.URL.Path)
			}
			w.WriteHeader(http.StatusNoContent)
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCommentTools(s, c)
		res := callTool(t, s, "clickup_delete_comment", map[string]any{"comment_id": "c1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if !strings.Contains(textOf(t, res), "c1") {
			t.Errorf("result = %q, want it to contain comment id", textOf(t, res))
		}
	})
}

func TestClickupListCommentReplies(t *testing.T) {
	t.Run("missing_required_arg", func(t *testing.T) {
		hit := false
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			hit = true
			t.Errorf("server should not have been called")
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCommentTools(s, c)
		res := callTool(t, s, "clickup_list_comment_replies", map[string]any{})
		if !res.IsError {
			t.Error("IsError = false, want true")
		}
		if hit {
			t.Error("server was hit despite missing comment_id")
		}
	})

	t.Run("wiring", func(t *testing.T) {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("method = %s, want GET", r.Method)
			}
			if r.URL.Path != "/comment/c1/reply" {
				t.Errorf("path = %s, want /comment/c1/reply", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"comments":[{"id":"r1"}]}`))
		})
		s := server.NewMCPServer("test", "1.0.0")
		RegisterCommentTools(s, c)
		res := callTool(t, s, "clickup_list_comment_replies", map[string]any{"comment_id": "c1"})
		if res.IsError {
			t.Fatalf("IsError = true, want false: %s", textOf(t, res))
		}
		if !strings.Contains(textOf(t, res), "r1") {
			t.Errorf("result = %q, want it to contain r1", textOf(t, res))
		}
	})
}

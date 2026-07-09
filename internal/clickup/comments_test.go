package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestListTaskCommentsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"comments":[{"id":"c1"}]}`))
	})

	out, err := c.ListTaskComments(context.Background(), "task1")
	if err != nil {
		t.Fatalf("ListTaskComments: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/task/task1/comment" {
		t.Errorf("path = %q, want /task/task1/comment", gotPath)
	}
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("out is not a map: %T", out)
	}
	if _, ok := m["comments"]; !ok {
		t.Errorf("decoded = %+v, missing comments", m)
	}
}

func TestCreateTaskCommentPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"c1"}`))
	})

	_, err := c.CreateTaskComment(context.Background(), "task1", map[string]any{"comment_text": "hello"})
	if err != nil {
		t.Fatalf("CreateTaskComment: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/task/task1/comment" {
		t.Errorf("path = %q, want /task/task1/comment", gotPath)
	}
	if gotBody["comment_text"] != "hello" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestListListCommentsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"comments":[]}`))
	})

	if _, err := c.ListListComments(context.Background(), "list1"); err != nil {
		t.Fatalf("ListListComments: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/list/list1/comment" {
		t.Errorf("path = %q, want /list/list1/comment", gotPath)
	}
}

func TestCreateListCommentPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"c1"}`))
	})

	_, err := c.CreateListComment(context.Background(), "list1", map[string]any{"comment_text": "note"})
	if err != nil {
		t.Fatalf("CreateListComment: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/list/list1/comment" {
		t.Errorf("path = %q, want /list/list1/comment", gotPath)
	}
	if gotBody["comment_text"] != "note" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestListViewCommentsHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"comments":[]}`))
	})

	if _, err := c.ListViewComments(context.Background(), "view1"); err != nil {
		t.Fatalf("ListViewComments: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/view/view1/comment" {
		t.Errorf("path = %q, want /view/view1/comment", gotPath)
	}
}

func TestCreateViewCommentPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"c1"}`))
	})

	_, err := c.CreateViewComment(context.Background(), "view1", map[string]any{"comment_text": "chat message"})
	if err != nil {
		t.Fatalf("CreateViewComment: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/view/view1/comment" {
		t.Errorf("path = %q, want /view/view1/comment", gotPath)
	}
	if gotBody["comment_text"] != "chat message" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestUpdateCommentPutsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"c1","resolved":true}`))
	})

	_, err := c.UpdateComment(context.Background(), "c1", map[string]any{"resolved": true})
	if err != nil {
		t.Fatalf("UpdateComment: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if gotPath != "/comment/c1" {
		t.Errorf("path = %q, want /comment/c1", gotPath)
	}
	if gotBody["resolved"] != true {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestDeleteCommentHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	if err := c.DeleteComment(context.Background(), "c1"); err != nil {
		t.Fatalf("DeleteComment: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/comment/c1" {
		t.Errorf("path = %q, want /comment/c1", gotPath)
	}
}

func TestListCommentRepliesHitsExpectedPath(t *testing.T) {
	var gotMethod, gotPath string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"comments":[]}`))
	})

	if _, err := c.ListCommentReplies(context.Background(), "c1"); err != nil {
		t.Fatalf("ListCommentReplies: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/comment/c1/reply" {
		t.Errorf("path = %q, want /comment/c1/reply", gotPath)
	}
}

func TestCreateCommentReplyPostsExpectedBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_, _ = w.Write([]byte(`{"id":"reply1"}`))
	})

	_, err := c.CreateCommentReply(context.Background(), "c1", map[string]any{"comment_text": "reply text"})
	if err != nil {
		t.Fatalf("CreateCommentReply: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/comment/c1/reply" {
		t.Errorf("path = %q, want /comment/c1/reply", gotPath)
	}
	if gotBody["comment_text"] != "reply text" {
		t.Errorf("body = %+v", gotBody)
	}
}

func TestCommentsMethodsReturnAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"err":"Comment not found","ECODE":"COMM_001"}`))
	})

	_, err := c.ListTaskComments(context.Background(), "task1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusNotFound || apiErr.ECode != "COMM_001" || apiErr.Err != "Comment not found" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
}

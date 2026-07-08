package clickup

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// CreateTaskAttachment uploads a local file as an attachment on a task.
// filePath must be readable by the server process (e.g. inside the
// container's /workspace tree). Uses multipart upload, rebuilt fresh on
// every retry attempt since multipart bodies aren't replayable.
// POST /task/{task_id}/attachment
func (c *Client) CreateTaskAttachment(ctx context.Context, taskID, filePath string) (any, error) {
	var out any
	err := c.doRaw(ctx, func(ctx context.Context) (*http.Request, error) {
		f, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", filePath, err)
		}
		defer f.Close()

		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)
		part, err := w.CreateFormFile("attachment", filepath.Base(filePath))
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(part, f); err != nil {
			return nil, err
		}
		if err := w.Close(); err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURLv2+"/task/"+taskID+"/attachment", &buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		return req, nil
	}, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

package clickup

import (
	"net/url"
	"strconv"
)

// addParam sets key=value if value is non-empty.
func addParam(q url.Values, key, value string) {
	if value != "" {
		q.Set(key, value)
	}
}

// addBoolParam sets key to value only when present is true (i.e. the caller
// explicitly supplied the param, as opposed to it defaulting to false).
func addBoolParam(q url.Values, key string, value, present bool) {
	if present {
		q.Set(key, strconv.FormatBool(value))
	}
}

// addIntParam sets key only when present is true.
func addIntParam(q url.Values, key string, value int, present bool) {
	if present {
		q.Set(key, strconv.Itoa(value))
	}
}

// addArrayParam appends one query entry per value under key, matching
// ClickUp's repeated-key array convention (e.g. "assignees[]").
func addArrayParam(q url.Values, key string, values []string) {
	for _, v := range values {
		if v != "" {
			q.Add(key, v)
		}
	}
}

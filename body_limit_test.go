package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleDecideBodyLimit(t *testing.T) {
	payload := []byte(`{"request":{"point":"DECISION_POINT_FLOOR","policy":{"floor":"1.0000"}}}`)
	for _, tc := range []struct {
		name                 string
		size, want           int
		unknownContentLength bool
	}{
		{"exact cap", maxBody, http.StatusOK, false},
		{"one byte over", maxBody + 1, http.StatusRequestEntityTooLarge, false},
		{"exact cap without content length", maxBody, http.StatusOK, true},
		{"one byte over without content length", maxBody + 1, http.StatusRequestEntityTooLarge, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := append(append([]byte{}, payload...), bytes.Repeat([]byte(" "), tc.size-len(payload))...)
			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/decide", bytes.NewReader(body))
			if tc.unknownContentLength {
				request.ContentLength = -1
			}
			handleDecide(response, request)
			if response.Code != tc.want {
				t.Fatalf("status=%d, want %d", response.Code, tc.want)
			}
		})
	}
}

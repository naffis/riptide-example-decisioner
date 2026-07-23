// Command decisioner-http is a public EXTERNAL_ENDPOINT Decisioner sample.
//
//	go run . -addr 127.0.0.1:8089
//
// Register with Riptide model config: {"url":"http://127.0.0.1:8089/decide"}
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const maxBody = 1 << 20

func main() {
	addr := flag.String("addr", "127.0.0.1:8089", "listen address")
	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("listen %s: %v", *addr, err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "LISTENING %s\n", lis.Addr())

	mux := http.NewServeMux()
	mux.HandleFunc("/decide", handleDecide)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func handleDecide(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBody))
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}
	var req DecideRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "malformed DecideRequest: "+err.Error(), http.StatusBadRequest)
		return
	}
	out, err := DecideFloor(r.Context(), req.Request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(DecideResponse{Response: out})
	if err != nil {
		http.Error(w, "marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

package httpserve

import (
	"io"
	"net/http"

	"git.xiaojukeji.com/pearls/tlog"
	"git.xiaojukeji.com/pearls/tlog/iface"
)

func init() {
	http.HandleFunc("/tlog/level", logLevel)
}

func logLevel(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getLogLevel(w, r)
	case http.MethodPost, http.MethodPut:
		setLogLevel(w, r)
	default:
		writeTextResponse(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func getLogLevel(w http.ResponseWriter, r *http.Request) {
	level := "unknown"
	status := http.StatusOK
	if gs, ok := tlog.GetLogger().(iface.GetSetLevel); ok {
		level = gs.GetLevel().String()
		status = http.StatusInternalServerError
	}
	writeTextResponse(w, status, level)
}

func setLogLevel(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		writeTextResponse(w, http.StatusInternalServerError, "Unexpected Request Body")
		return
	}

	level, err := iface.StringToLevel(string(body))
	if err != nil {
		writeTextResponse(w, http.StatusInternalServerError, "Unknown Log Level")
		return
	}

	logger := tlog.GetLogger()
	if gs, ok := logger.(iface.GetSetLevel); ok {
		gs.SetLevel(level)
		writeTextResponse(w, http.StatusOK, "OK")
		return
	}
	writeTextResponse(w, http.StatusInternalServerError, "Unexpected Logger")
}

func writeTextResponse(w http.ResponseWriter, status int, txt string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(txt))
}

package tlog

import (
	"io"
	"net/http"

	"git.xiaojukeji.com/pearls/tlog/iface"
)

func init() {
	http.HandleFunc("/tlog/level/get", GetLogLevel)
	http.HandleFunc("/tlog/level/set", SetLogLevel)
}

func GetLogLevel(w http.ResponseWriter, r *http.Request) {
	level := "unknown"
	status := http.StatusOK
	if gs, ok := GetLogger().(iface.GetSetLevel); ok {
		level = gs.GetLevel().String()
		status = http.StatusInternalServerError
	}
	writeTextResponse(w, status, level)
}

func SetLogLevel(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeTextResponse(w, http.StatusInternalServerError, "Unexpected Request Body")
		return
	}

	level, err := iface.StringToLevel(string(body))
	if err != nil {
		writeTextResponse(w, http.StatusInternalServerError, "Unknown Log Level")
		return
	}

	logger := GetLogger()
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

package app

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed client/*
var client embed.FS
var contentFS, _ = fs.Sub(client, "client")

func (app *App) HandleIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Length", "12")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("Hello World!"))
}

func (app *App) HandleIndex2(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("Hello 2 tesg!"))
}

func HandleClient() http.Handler {
	fileServer := http.FileServer(http.FS(contentFS))
	return http.StripPrefix("/", fileServer)

	//http.FileServer(http.Dir("./client"))
}

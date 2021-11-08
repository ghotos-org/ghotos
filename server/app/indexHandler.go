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

func (frw *FallbackResponseWriter) WriteHeader(statusCode int) {
	//log.Printf("INFO: WriteHeader called with code %d\n", statusCode)
	if statusCode == http.StatusNotFound {
		//log.Printf("INFO: Setting FileNotFound flag\n")
		frw.FileNotFound = true
		return
	}
	frw.WrappedResponseWriter.WriteHeader(statusCode)
}

type (
	// FallbackResponseWriter wraps an http.Requesthandler and surpresses
	// a 404 status code. In such case a given local file will be served.
	FallbackResponseWriter struct {
		WrappedResponseWriter http.ResponseWriter
		FileNotFound          bool
	}
)

// Header returns the header of the wrapped response writer
func (frw *FallbackResponseWriter) Header() http.Header {
	return frw.WrappedResponseWriter.Header()
}

// Write sends bytes to wrapped response writer, in case of FileNotFound
// It surpresses further writes (concealing the fact though)
func (frw *FallbackResponseWriter) Write(b []byte) (int, error) {
	if frw.FileNotFound {
		return len(b), nil
	}
	return frw.WrappedResponseWriter.Write(b)
}

func (app *App) HandleClient(w http.ResponseWriter, r *http.Request) {

	frw := FallbackResponseWriter{
		WrappedResponseWriter: w,
		FileNotFound:          false,
	}

	http.FileServer(http.FS(contentFS)).ServeHTTP(&frw, r)
	if frw.FileNotFound {
		b, _ := client.ReadFile("client" + "/index.html")
		w.Header().Set("Content-Type", "text/html")
		w.Write(b)
		return
	}

}

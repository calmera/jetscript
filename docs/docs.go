package docs

import (
	"embed"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"io/fs"
	"net/http"
	"os"
	"time"
)

//go:embed web
var web embed.FS

func ServeDocs(port int, docsDir string) error {
	// check if the docs directory exists
	_, err := fs.Stat(web, docsDir)
	if err != nil {
		return fmt.Errorf("failed to stat docs dir: %w", err)
	}

	// check if the index.json file exists within the docs dir
	_, err = fs.Stat(web, fmt.Sprintf("%s/index.json", docsDir))
	if err != nil {
		return fmt.Errorf("failed to stat docs dir: missing index: %w", err)
	}

	address := fmt.Sprintf(":%d", port)

	server := &http.Server{
		Addr:         address,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 0,
		IdleTimeout:  10 * time.Second,
	}

	router := mux.NewRouter()

	dist, err := fs.Sub(web, "web/dist")
	if err != nil {
		return fmt.Errorf("failed to navigate web fs: %w", err)
	}

	docsFs := os.DirFS(docsDir)

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.FS(dist))))
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.FS(docsFs))))
	//router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//  http.ServeFile(w, r, "./build/index.html")
	//})

	log.Info().Msgf("api running on port %d", port)
	server.Handler = router

	return server.ListenAndServe()
}

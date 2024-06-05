package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

var staticDir string
var serverPort string

func main() {

	staticDir := "web/dist"
	
	router := chi.NewRouter()
	routes := Routes{
		staticDir: staticDir,
		apiKey: "key",
		apiSecret: "secret",
	}

	b, err := os.ReadFile("client_secret.json")
	if err != nil {
	  log.Fatalf("Unable to read client secret file: %v", err)
	}
  
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
	  log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	router.Get("/", func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, staticDir + "/index.html")
	})

	router.Get("/auth", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		auth_url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		// sanitize URL
		parsed_url, err := url.Parse(auth_url)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(resp, req, parsed_url.String(), http.StatusFound)
	}))

	router.Get("/", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		queryParams := req.URL.Query()
	
		code := queryParams.Get("code")
	
		if code == "" {
			http.Error(resp, "Missing 'code' query parameter", http.StatusBadRequest)
			return
		}
		tok := saveCodeAsToken(config, code)
		cacheFile, err := tokenCacheFile()
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		saveToken(cacheFile, tok)
		// Write the extracted code to the response
		http.Redirect(resp, req, "http://localhost:5173", http.StatusFound)
	}))

	router.Get("/api/subscriptions", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := getCachedClient(r.Context(), config)
		service, err := youtube.New(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		call := service.Subscriptions.List([]string{"id", "snippet", "contentDetails"}).Mine(true)
		response, err := call.Do()
		if err != nil {
		  log.Fatalf("Error calling API: %v", err)
		}
		json.NewEncoder(w).Encode(response)
	}))

	router.Get("/api/playlists", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := getCachedClient(r.Context(), config)
		service, err := youtube.New(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		call := service.Playlists.List([]string{"id", "snippet", "contentDetails"}).Mine(true)
		response, err := call.Do()
		if err != nil {
		  log.Fatalf("Error calling API: %v", err)
		}
		json.NewEncoder(w).Encode(response)
	}))

	// Serve static files
	FileServer(router, "/", http.Dir(staticDir))

	// Start server
	serverPort = os.Getenv("PORT")
	if len(serverPort) == 0 {
		serverPort = "3333"
	}

	fmt.Printf("### Starting server listening on %v\n", serverPort)
	fmt.Printf("### Serving static content from '%v'\n", staticDir)
	fmt.Printf("### Browse: http://localhost:3333\n")
	http.ListenAndServe(":"+serverPort, router)
}

// Serve static files
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		fs.ServeHTTP(resp, req)
	}))

}

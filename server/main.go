package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

var serverPort string

// allow vue to access server data
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

		// Proceed to the next middleware or handler
		next.ServeHTTP(w, r)
	})
}

func main() {

	staticDir := "web/dist"

	router := chi.NewRouter()
	// routes := Routes{
	// 	staticDir: staticDir,
	// 	apiKey: "key",
	// 	apiSecret: "secret",
	// }

	router.Use(middleware.Logger)

	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		slog.Warn("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	router.Get("/", func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, staticDir+"/index.html")
	})

	router.Get("/api/auth", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		auth_url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		// sanitize URL
		parsed_url, err := url.Parse(auth_url)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(resp, req, parsed_url.String(), http.StatusFound)
	}))

	router.Get("/api/auth/profile", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		youtubeService, err := generateService(config, resp)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}

		call := youtubeService.Channels.List([]string{"id", "snippet"}).Mine(true)
		response, err := call.Do()
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(resp).Encode(response)

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
		call := service.Subscriptions.List([]string{"id", "snippet", "contentDetails"}).Mine(true).MaxResults(50)
		response, err := call.Do()
		if err != nil {
			slog.Error("Error calling Subscriptions API: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var subscriptions []Subscription
		for _, item := range response.Items {
			if item == nil || item.Snippet == nil || item.Snippet.Thumbnails.Default == nil {
				slog.Error("item or item.SubscriberSnippet or item.SubscriberSnippet.Thumbnails.Default is nil")
				continue
			}

			subscriptions = append(subscriptions, Subscription{
				Title:         item.Snippet.Title,
				ThumbnailLink: item.Snippet.Thumbnails.Default.Url,
			})
		}

		json.NewEncoder(w).Encode(subscriptions)
	}))

	// router.Post("/api/subscriptions/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	client := getCachedClient(r.Context(), config)
	// 	service, err := youtube.New(client)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	}
	// 	id := chi.URLParam(r, "id")
	// 	subscription := youtube.Subscription{
	// 		Id: id,
	// 	}
	// 	call := service.Subscriptions.Insert([]string{"id", "snippet", "contentDetails"}, &subscription)
	// 	response, err := call.Do()
	// 	if err != nil {
	// 		slog.Error("Error calling Subscription Insertion API: %v", err)
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	json.NewEncoder(w).Encode(response)
	// }))

	router.Get("/api/playlists", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		youtubeService, err := generateService(config, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		call := youtubeService.Playlists.List([]string{"id", "snippet", "contentDetails"}).Mine(true).MaxResults(50)
		raw_playlists, err := call.Do()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error calling API: %v", err), http.StatusInternalServerError)
			return
		}

		if raw_playlists == nil || raw_playlists.Items == nil {
			http.Error(w, "No playlists found", http.StatusNotFound)
			return
		}

		var allPlaylists []Playlist

		for _, item := range raw_playlists.Items {
			if item == nil || item.Snippet == nil || item.Snippet.Thumbnails == nil || item.Snippet.Thumbnails.Default == nil {
				continue
			}

			raw_videos := youtubeService.PlaylistItems.List([]string{"id", "snippet"}).PlaylistId(item.Id).MaxResults(50)
			videosResponse, err := raw_videos.Do()
			if err != nil {
				http.Error(w, fmt.Sprintf("Error calling API: %v", err), http.StatusInternalServerError)
				return
			}

			if videosResponse == nil || videosResponse.Items == nil {

				continue
			}

			var allVideos []Video

			fmt.Println("length of videosResponse.Items is ", len(videosResponse.Items))
			for _, videoItem := range videosResponse.Items {
				if videoItem == nil || videoItem.Snippet == nil || videoItem.Snippet.ResourceId == nil || videoItem.Snippet.Thumbnails == nil {
					slog.Warn("Skipping video")
					fmt.Println("Skipping video")

					continue
				}

				var link string
				if videoItem.Snippet.Thumbnails.High != nil {
					link = videoItem.Snippet.Thumbnails.High.Url
				}
				if link == "" && videoItem.Snippet.Thumbnails.Default != nil {
					link = videoItem.Snippet.Thumbnails.Default.Url
				}

				video := Video{
					ID:            videoItem.Snippet.ResourceId.VideoId,
					Title:         videoItem.Snippet.Title,
					Link:          "https://www.youtube.com/watch?v=" + videoItem.Snippet.ResourceId.VideoId,
					ThumbnailLink: link,
				}
				allVideos = append(allVideos, video)
			}

			playlist := Playlist{
				Videos:        allVideos,
				ThumbnailLink: item.Snippet.Thumbnails.High.Url,
				Title:         item.Snippet.Title,
			}

			allPlaylists = append(allPlaylists, playlist)
		}

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "\t")
		if err := encoder.Encode(allPlaylists); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
			return
		}
	}))

	router.Get("/api/playlists/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	}))

	router.Get("/api/search", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, "Query parameter is required", http.StatusBadRequest)
			return
		}

		youtubeService, err := generateService(nil, w) // Pass appropriate config instead of nil
		if err != nil {
			slog.Error("Error generating YouTube service: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		call := youtubeService.Search.List([]string{"id", "snippet"})
		if call == nil {
			slog.Error("Error creating API call: call is nil")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		call = call.Q(query).MaxResults(50)
		if call == nil {
			slog.Error("Error setting query or max results: call is nil")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response, err := call.Do()
		if err != nil {
			slog.Error("Error calling YouTube API: %v", err)
			http.Error(w, fmt.Sprintf("Error calling API: %v", err), http.StatusInternalServerError)
			return
		}

		if response == nil || response.Items == nil {
			slog.Error("No results found")
			http.Error(w, "No results found", http.StatusNotFound)
			return
		}

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "\t")
		if err := encoder.Encode(response.Items); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
			return
		}
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
	fmt.Printf("### Browse: http://localhost:%v\n", serverPort)
	http.ListenAndServe(":"+serverPort, corsMiddleware(router))
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

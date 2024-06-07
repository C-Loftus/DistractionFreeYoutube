// Sample Go code for user authorization

package main

import (
	"encoding/json"
	"fmt"
	"log"

	// "net/http"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

func generateService(config *oauth2.Config, w http.ResponseWriter) (*youtube.Service, error) {
	ctx := context.Background()
	token := getCachedToken(ctx, config)
	if token == nil {
		http.Error(w, "Error getting cached token", http.StatusInternalServerError)
		err := fmt.Errorf("error getting cached token")
		return nil, err
	}

	youtubeService, err := youtube.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		err := fmt.Errorf("error creating YouTube client: %v", err)
		return nil, err
	}

	return youtubeService, nil

}

func getCachedClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, _ := tokenFromFile(cacheFile)
	return config.Client(ctx, tok)
}

func getCachedToken(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, _ := tokenFromFile(cacheFile)
	return tok
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func saveCodeAsToken(config *oauth2.Config, code string) *oauth2.Token {

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to convert code to oauth token: %v", err)
	}
	return tok
}

func tokenCacheFile() (string, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// Define the token cache directory within the current working directory
	tokenCacheDir := filepath.Join(cwd, ".credentials")
	// Create the token cache directory with appropriate permissions
	err = os.MkdirAll(tokenCacheDir, 0700)
	if err != nil {
		return "", err
	}
	// Return the path to the token cache file
	return filepath.Join(tokenCacheDir, url.QueryEscape("youtube-credentials.json")), nil
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

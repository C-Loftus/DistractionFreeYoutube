package main

type Video struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Link    string `json:"link"`
	ThumbnailLink string `json:"thumbnailLink"`
}

type Playlist struct {
	Title string `json:"title"`
	Videos []Video `json:"videos"`
	ThumbnailLink string `json:"thumbnailLink"`
}

type Subscription struct {
	Title string `json:"title"`
	ThumbnailLink string `json:"thumbnailLink"`
	
}
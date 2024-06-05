package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Routes struct {
	staticDir string
	apiKey    string
	apiSecret string
}

// Route: Data
type Data struct {
	Query    string `json:"query"`
	Datetime string `json:"datetime"`
}

// /api/data
func (r Routes) apiDataRoute(resp http.ResponseWriter, req *http.Request) {
	var data Data

	resp.Header().Set("Access-Control-Allow-Origin", "*")

	data.Query = "Dummy data"
	datetime := time.Now().Format("Mon Aug 5 10:17:08 EST 2019")
	data.Datetime = datetime

	// Reformat data as JSON
	result, err := json.Marshal(data)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	// print result
	fmt.Println(string(result))

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(result)
}

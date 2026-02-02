package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type MovieData struct {
	ID int `json:"id"`
}

type RequestBody struct {
	EventType    *string    `json:"eventType"`
	InstanceName *string    `json:"instanceName"`
	Movie        *MovieData `json:"movie"`
	DownloadId   *string     `json:"downloadId"`
	DeleteFiles  bool       `json:"deletedFiles"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	parsedBody := RequestBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parsedBody)
	if err != nil || parsedBody.EventType == nil || parsedBody.Movie == nil || parsedBody.InstanceName == nil {
		fmt.Printf("JSON validation error: %v\n", err)
		w.WriteHeader(400)
		return
	}

	if *parsedBody.EventType == "Grab" && parsedBody.DownloadId != nil {
		db.AddMovie(MapData{
			Instance:    *parsedBody.InstanceName,
			InstanceId:  parsedBody.Movie.ID,
			TorrentHash: strings.ToLower(*parsedBody.DownloadId),
		})
	} else if *parsedBody.EventType == "MovieDelete" {
		db.RemoveMovie(*parsedBody.InstanceName, parsedBody.Movie.ID, parsedBody.DeleteFiles)
	}
	w.WriteHeader(200)
}

func main() {
	db.Load()
	http.HandleFunc("/", handleRequest)
	err := http.ListenAndServe(":9191", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

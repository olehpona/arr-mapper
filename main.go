package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type MovieData struct {
	ID int `json:"id"`
}

type RequestBody struct {
	EventType    *string    `json:"eventType"`
	InstanceName *string    `json:"instanceName"`
	Movie        *MovieData `json:"movie"`
	DownloadId   *string    `json:"downloadId"`
	DeleteFiles  bool       `json:"deletedFiles"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	parsedBody := RequestBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parsedBody)
	if err != nil || parsedBody.EventType == nil || parsedBody.Movie == nil || parsedBody.InstanceName == nil {
		logrus.Errorf("JSON validation error: %v", err)
		w.WriteHeader(400)
		return
	}

	if *parsedBody.EventType == "Grab" && parsedBody.DownloadId != nil {
		logrus.Infof("Adding movie %d from instance %s with torrent hash %s", parsedBody.Movie.ID, *parsedBody.InstanceName, *parsedBody.DownloadId)
		db.AddMovie(MapData{
			Instance:    *parsedBody.InstanceName,
			InstanceId:  parsedBody.Movie.ID,
			TorrentHash: strings.ToLower(*parsedBody.DownloadId),
		})
	} else if *parsedBody.EventType == "MovieDelete" {
		logrus.Infof("Removing movie %d from instance %s, delete files: %t", parsedBody.Movie.ID, *parsedBody.InstanceName, parsedBody.DeleteFiles)
		db.RemoveMovie(*parsedBody.InstanceName, parsedBody.Movie.ID, parsedBody.DeleteFiles)
	}
	w.WriteHeader(200)
}

func initLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.WarnLevel)
}

func main() {
	initLogger()
	err := db.Load()
	if err != nil {
		logrus.Errorf("error loading database: %s", err)
	}
	http.HandleFunc("/", handleRequest)
	logrus.Warnf("Server running on port 9191")
	err = http.ListenAndServe(":9191", nil)
	if errors.Is(err, http.ErrServerClosed) {
		logrus.Infof("server closed")
	} else if err != nil {
		logrus.Errorf("error starting server: %s", err)
		os.Exit(1)
	}
}

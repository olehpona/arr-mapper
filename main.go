package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"github.com/sirupsen/logrus"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	parsedBody := RequestBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parsedBody)
	if err != nil || parsedBody.EventType == nil || parsedBody.InstanceName == nil {
		logrus.Errorf("JSON validation error: %v", err)
		w.WriteHeader(400)
		return
	}
	
	var instanceId int

	if parsedBody.Movie != nil {
		instanceId = parsedBody.Movie.ID
	} else if parsedBody.Series != nil {
		instanceId = parsedBody.Series.ID
	} else {
		logrus.Errorf("No movie or series data provided")
		w.WriteHeader(400)
		return
	}

	if *parsedBody.EventType == "Grab" && parsedBody.DownloadId != nil {
		logrus.Infof("Adding movie %d from instance %s with torrent hash %s", instanceId, *parsedBody.InstanceName, *parsedBody.DownloadId)
		db.AddMovie(MapData{
			Instance: *parsedBody.InstanceName,
			InstanceId:       instanceId,
			TorrentHash: strings.ToLower(*parsedBody.DownloadId),
		})
	} else if *parsedBody.EventType == "MovieDelete" || *parsedBody.EventType == "SeriesDelete" {
		logrus.Infof("Removing movie %d from instance %s, delete files: %t", instanceId, *parsedBody.InstanceName, parsedBody.DeleteFiles)
		db.RemoveMovie(*parsedBody.InstanceName, instanceId, parsedBody.DeleteFiles)
	}
	w.WriteHeader(200)
}

func initLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
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

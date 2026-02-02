package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func removeTorrent(hash string, withData bool) error {
	url := os.Getenv("TRANSMISSION_URL")
    if url == "" {
        url = "http://localhost:9091/transmission/rpc"
    }
	
	payload, _ := json.Marshal(map[string]interface{}{
		"method": "torrent-remove",
		"arguments": map[string]interface{}{
			"ids":               []string{hash},
			"delete-local-data": withData,
		},
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusConflict {
		sessionID := resp.Header.Get("X-Transmission-Session-Id")
		resp.Body.Close()
		req, _ = http.NewRequest("POST", url, bytes.NewBuffer(payload))
		req.Header.Set("X-Transmission-Session-Id", sessionID)
		resp, err = client.Do(req)
	}

	if resp != nil {
		defer resp.Body.Close()
	}
	return err
}
package main

import (
	"encoding/json"
	"os"
	"sync"
)

type MapData struct {
	Instance string    `json:"instance"`
	InstanceId int     `json:"instance_id"`
	TorrentHash string `json:"torrent_hash"`
}

type Database struct {
	mu     sync.Mutex
	Map []MapData `json:"movies"`
}

var db = &Database{Map: make([]MapData, 0)}
const dbFile = "data.json"

func (d *Database) Load() error {
	file, err := os.ReadFile(dbFile)
	if err != nil {
		if os.IsNotExist(err) { return nil }
		return err
	}
	return json.Unmarshal(file, d)
}

func (d *Database) Save() error {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil { return err }
	return os.WriteFile(dbFile, data, 0644)
}


func (d *Database) AddMovie(m MapData) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	d.Map = append(d.Map, m)
	return d.Save()
}

func (d *Database) RemoveMovie(instance string, instanceId int, withData bool) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	index := -1
	var hashToDelete string

	for i, v := range d.Map {
		if v.Instance == instance && v.InstanceId == instanceId {
			index = i
			hashToDelete = v.TorrentHash
			break
		}
	}

	if index != -1 {
		if hashToDelete != "" {
			go removeTorrent(hashToDelete, withData)
		}

		d.Map = append(d.Map[:index], d.Map[index+1:]...)
		return d.Save()
	}

	return nil
}
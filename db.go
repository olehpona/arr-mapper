package main

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"
)

type MapData struct {
	Instance    string `json:"instance"`
	InstanceId  int    `json:"instance_id"`
	TorrentHash string `json:"torrent_hash"`
}

type Database struct {
	mu  sync.Mutex
	Map map[string][]string `json:"movies"`
}

var db = &Database{Map: make(map[string][]string)}

const dbFile = "data.json"

func (d *Database) Load() error {
	file, err := os.ReadFile(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	err = json.Unmarshal(file, d)
	if d.Map == nil {
		d.Map = make(map[string][]string)
	}
	return err
}

func (d *Database) Save() error {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dbFile, data, 0644)
}

func (d *Database) createKey(instance string, instanceId int) string {
	return instance + ":" + strconv.Itoa(instanceId)
}

func (d *Database) findIndexLocked(key string, hash string) int {
	for i, v := range d.Map[key] {
		if v == hash {
			return i
		}
	}
	return -1
}

func (d *Database) AddMovie(m MapData) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	key := d.createKey(m.Instance, m.InstanceId)
	idx := d.findIndexLocked(key, m.TorrentHash)
	if idx != -1 {
		return nil
	}

	d.Map[key] = append(d.Map[key], m.TorrentHash)
	return d.Save()
}

func (d *Database) RemoveMovie(instance string, instanceId int, withData bool) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	key := d.createKey(instance, instanceId)

	for _, hash := range d.Map[key] {
		go removeTorrent(hash, withData)
	}
	delete(d.Map, key)
	return d.Save()
}

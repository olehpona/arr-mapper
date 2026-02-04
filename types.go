package main

type MovieData struct {
	ID int `json:"id"`
}

type SeriesData struct {
	ID int `json:"id"`
}

type RequestBody struct {
	EventType    *string     `json:"eventType"`
	InstanceName *string     `json:"instanceName"`
	Movie        *MovieData  `json:"movie"`
	Series 	     *SeriesData `json:"series"`
	DownloadId   *string     `json:"downloadId"`
	DeleteFiles  bool        `json:"deletedFiles"`
}
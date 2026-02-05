![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
# ARR-MAPPER
arr mapper - is simple http server created to map arr family apps ( radarr, sonarr ) media id with coresponding torrent hashes.

## WHY ?
It is used for posibility to remove movie from radarr ant torrent client ( currently only transmission ) with only one click in radarr/sonnar

## Features
- Radarr, Sonarr support
- backing up state in json file to prevent data lost
- Transmission support

## How to use
- Set env variable
```sh
export TRANSMISSION_URL="http://localhost:9091/transmission/rpc"
```
- Run binarry 
```sh
./arr-mapper
```
The server runs on port 9191 by default.
- Configure Webhooks: In Radarr/Sonarr, go to Settings - Connect - Add Webhook:

URL: http://arr-mapper-path:9191/  
METHOD: POST
## Required webhook events
### Radarr
- OnGrab
- OnMovieDelete
### Sonarr
- OnGrab
- OnSeriesDelete
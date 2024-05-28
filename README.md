# GoCam

Program to record cameras via RSTP with ffmpeg

## Endpoints

- `/`
  - Index page
- `/listcameras`
  - List all available cameras from config file
- `/listrecordings`
  - List current recordings
- `/record`
  - Start recording a specific camera for a specific time
- `/cancel`
  - Cancel a current recording
  
## Usage

- Define a toml config file

``` toml
[conf]
json_file = "/tmp/file.json"
listen_port = "9999"
dump_recording = "/tmp/dump/"
log_recording = "/tmp/log/"

[[targets]]
camera_id = "1"
name = "Bedroom Tapo Camera 1"
host = "192.168.30.44"
port = "554"
stream = "/stream1"
protocol = "rtsp"
user = "brun0teixeira"
password = "qwerty654321"
recording_path = "/mnt/external/camera/bedroom-tapo/"

[[targets]]
camera_id = "2"
name = "Bedroom Tapo Camera 2"
host = "192.168.30.45"
port = "554"
stream = "/stream1"
protocol = "rtsp"
user = "brun0teixeira"
password = "1qazZAQ!"
recording_path = "/mnt/external/camera/bedroom-tapo/"
```

- Run the program

``` console
$ gocam -f config.toml
```

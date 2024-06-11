# GoCam

Program to record cameras via RTSP with ffmpeg

![Untitled-2024-06-11-1602](https://github.com/BrunoTeixeira1996/gocam/assets/12052283/1a85532c-0d32-4704-8ec0-8a08a937ed72)


## Endpoints

- `/`
  - Index page
- `/listcameras`
  - List all available cameras from config file
  - ![image](https://github.com/BrunoTeixeira1996/gocam/assets/12052283/db129ac8-155c-49d0-a905-eba9a5d352c2)
- `/listrecordings`
  - List current recordings
  - ![image](https://github.com/BrunoTeixeira1996/gocam/assets/12052283/1c78957e-5e1b-4fe7-81f3-4412715c36dd)
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

- Start recording

```console
$ curl http://<IP>:9999/record/ -X POST -d '{"cameraId": "1", "duration":"2m"}' -m 1
```

- Stop recording

```console
$ curl http://<IP>:9999/cancel/ -X POST -d '{"ID":"sN3eYaFF8u"}'
```

# nas
small go nas platform for my raspberry pi

## usage

```
cp assets/config/config.example.json assets/config/config.json
 # edit config file with your hot/cold storage locations
nano assets/config/config.json
# run
go run webserver.go
# or build and run
go build; ./nas
```

you can also optionally use the `build-pi.sh` to build it for a raspberry pi (tested with raspberry pi 3 model b+)

then navigate to `localhost:3000`

## credits

svg icons via https://iconsvg.xyz

raspberry pi svg via https://www.vectorlogo.zone/logos/raspberrypi/index.html
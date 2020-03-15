# nas
small go nas platform for my raspberry pi

## usage

### configuration

unlike the initial version of this project, the current build uses _providers_ to determine how to handle various 
functions related to files. currently, two are implemented, `disk` and `backblaze`, since they are the primary providers
i use myself. the providers you would like to use can be added to `providers.yml` alongside the binary.

for example, here is a sample configuration implementing both of them:

```yaml
disk:
  provider: disk
  path: /tmp/nas
backblaze:
  provider: backblaze
  config:
    appKeyId: APP_KEY_ID
    appId: APP_ID
    bucket: BUCKET_ID
```

(read more here: [#providers](#providers))

### running

after adding the providers you would like to use, the application can be run simply with `./nas`.

### building

this project uses go modules and a makefile, so building should be relatively straightforward. 

 - `make` will build the project for your system's architecture.
 - `make pi` will build the project with the `GOOS=linux GOARCH=arm GOARM=5 go` flags set for raspberry pis.

## api

initially the heavy lifting was done by the server, but the need for a better frontend was clear.

full documentation coming soon once actual functionality has been nailed down.

## providers

// todo

## credits

svg icons via https://iconsvg.xyz

raspberry pi svg via https://www.vectorlogo.zone/logos/raspberrypi/index.html

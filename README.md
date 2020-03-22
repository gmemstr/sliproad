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
 
### providers

there are a few built-in providers, and more can be added by opening a pull request.

|name|service|configuration example|
|----|-------|---------------------|
|disk|local filesystem|disk.yml|
|backblaze|backblaze b2|backblaze.yml|

#### custom provider

custom file providers can be implemented by adding a new go file to the `files` module. it should
implement the `FileProviderInterface` interface.
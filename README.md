# nas
bringing filesystems together

## about

this project aims to be an easy-to-manage web application that allows the management of cloud storage, whether it be on
the host machine or part of a remote api. this is intended mostly for my own use, but i am documenting it in a way that
i hope allows others to pick it up and improve on it down the line.

if something is unclear, feel free to open an issue :)

## configuration

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

## running

after adding the providers you would like to use, the application can be run simply with `./nas`. it will attach to port
`:3000`.

## building

this project uses go modules and a makefile, so building should be relatively straightforward. 

 - `make` will build the project for your system's architecture.
 - `make run` will run the project with `go run`
 - `make pi` will build the project with the `GOOS=linux GOARCH=arm GOARM=5 go` flags set for raspberry pis.
 
## providers

"providers" provide a handful of functions to interact nicely with a filesystem, whether it be on local disk or on a 
remote server via an api. best-effort is done to keep these performant, up to date and minimal.

there are a few built-in providers, and more can be added by opening a pull request.

|name|service|configuration example|
|----|-------|---------------------|
|disk|local filesystem|assets/config_examples/disk.yml|
|backblaze|backblaze b2|assets/config_examples/backblaze.yml|

you can find a full configuration file under `assets/config_examples/providers.yml`

### custom provider

custom file providers can be implemented by adding a new go file to the `files` module. it should
implement the `FileProviderInterface` interface.

## authentication

basic authentication support utilizing [keycloak](https://keycloak.org/) has been implemented, but work
is being done to bring this more inline with the storage provider implementation. see `assets/config_examples/auth.yml`
for an example configuration - having this file alongside the binary will activate authentication on all
`/api/files` endpoints. note that this implementation is a work in progress, and a seperate branch will
contain further improvements.

## icons

SVG Icons provided by Paweł Kuna: https://github.com/tabler/tabler-icons (see assets/web/icons/README.md)
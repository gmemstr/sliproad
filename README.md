# Sliproad

Merging filesystems together

## About

This project aims to be an easy-to-use web API and frontend that allows the
management of cloud storage alongside local filesystems. While this is intended 
mostly for my own use, I am documenting it in a way that I hope allows others 
to use it!

## Configuration

Sliproad uses "Providers" to support various filesystems "types", whether it be
remote or local. Currently, three exist - `disk` for filesystems local to the
machine, `backblaze` to leverage Backblaze B2 file storage and `s3` for AWS S3
(and other compatible providers).

An example of leveraging all three, in various forms, can be found below. As 
more are added, this example will be updated, and more examples can be found in
the `assets/config_examples` directory.

```yaml
disk:
  provider: disk
  path: /tmp/nas
backblaze:
  provider: backblaze
  config:
    bucket: some-bucket
    applicationKeyId: application-key-id
    applicationKey: application-key
s3:
  provider: s3
  config:
    region: eu-west-2
    bucket: some-bucket
# An example of an S3 compatible API, doesn't have to be Backblaze.
backblazes3:
  provider: s3
  config:
    bucket: some-bucket
    region: us-west-000
    endpoint: s3.us-west-000.backblazeb2.com
    keyid: key-id
    keysecret: key-secret
```

## Running

After configuring the providers you would like to utilize, simply run 
`./sliproad`. This will spin up the webserver at `127.0.0.1:3000`, listening on
all addresses.

## Frontend

The frontend is a very lightweight JavaScript application and aims to be very
functional, if a bit rough around the edges.

![Screenshot_2021-05-29 Sliproad](https://user-images.githubusercontent.com/1878840/120085420-d63cbc80-c0cf-11eb-9fbb-b0b05a3f5d58.png)

It should scale reasonably well for smaller devices. Because it's now bundled
into the binary (as opposed to distributed alongside), it's no longer possible
to swap it out for a custom frontend without serving the frontend seperately.

## API

This project is largely API-first, and documentation can be found here:

https://github.com/gmemstr/sliproad/wiki/API

## Building

This project leverages a Makefile to macro common commands for running, testing
and building this project.

 - `make` will build the project for your system's architecture.
 - `make run` will run the project with `go run`
 - `make pi` will build the project with the `GOOS=linux GOARCH=arm GOARM=5 go` flags set for Raspberry Pi.
 - `make dist` will build and package the binaries for distribution.

### Adding Providers

New file providers can be implemented by building off the 
`FileProviderInterface` struct, as the existing providers demonstrate. You can
then instruct the [`TranslateProvider()`](https://github.com/gmemstr/sliproad/blob/master/files/fileutils.go#L8-L21) function
that it exists and how to configure it.

## Authentication [!]

Authentication is a bit tricky and due to be reworked in the next iteration of
this project. Currently, support for [Keycloak](https://www.keycloak.org/) is
implemented, if a bit naively. You can turn this authentication requirement on
by adding `auth.yml` alongside your `providers.yml` file with the following:

```yaml
provider_url: "https://url-of-keycloak"
realm: "keycloak-realm"
redirect_base_url: "https://location-of-sliproad"
```

Keycloak support is not currently actively supported, and is due to be removed 
in the next major release of Sliproad. That said, if you encounter any major 
bugs utilizing it before this, _please_ open an issue so I can dig in further.

## Credits

SVG Icons provided by Paweł Kuna: https://github.com/tabler/tabler-icons (see assets/web/icons/README.md)

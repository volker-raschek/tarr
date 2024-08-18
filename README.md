# tarr

[![Build Status](https://drone.cryptic.systems/api/badges/volker.raschek/tarr/status.svg)](https://drone.cryptic.systems/volker.raschek/tarr)
[![Docker Pulls](https://img.shields.io/docker/pulls/volkerraschek/tarr)](https://hub.docker.com/r/volkerraschek/tarr)

The tarr project contains small binaries / tools for interacting with *arr applications. The tools are helpful in a
kubernetes environment to retrofit missing functions of the \*arr applications.

> [!NOTE]
> Instead of compiling the tarr applications by yourself, use the tarr container image instead. More described [below](#container-image).

## autharr

The binary `autharr` is a small program to extract from a `config.xml` or `config.yaml` the API token. The token is
written to the standard output. Alternatively, it can also be written to a file.

With regard to [exportarr](https://github.com/onedr0p/exportarr), it can be helpful in the Kubernetes environment to
extract the token and use it for other API queries. For example for healthchecks. It therefore solve the following
[problem](https://github.com/onedr0p/exportarr/issues/294).

```bash
$ autharr /etc/bazarr/config.yaml
do7IuHiewooFaiyu
$ autharr /etc/lidarr/config.xml
aeteipei4Meing5i
```

Alternatively, the `--watch` flag can be set. This monitors the config file and writes the API token to the defined
output in the event of changes.

```bash
$ autharr --watch /etc/bazarr/config.yaml
baGohkie9EL5Tahr
oov1liQuaiki1lar
vaeGa9Cheeheev2I
```

Pipe the output direct into a file. Exit the program by Ctrl+C.

```bash
$ autharr --watch /etc/bazarr/config.yaml /tmp/bazarr/token
^C
$
```

## healarr

The binary `healarr` is a small program to check if the *arr application is healthy. Some \*arr applications does not
have implemented a dedicated REST endpoint for healthchecks or like the liveness or readiness probe. Instead will be
called the API for a status, which returns 200 if the \*arr instance is healthy.

`healarr` uses the internal packages from `autharr` to extract the API token from a config file. Alternatively can
directly passed the API token as flag.

```bash
$ if healarr bazarr https://bazarr.example.com --config /etc/bazarr/config.xml; then
>   echo "Healthy"
> else
>   echo "Unhealthy"
> fi
Healthy
```

## container-image

The container image `docker.io/volkerraschek/tarr` contains all tarr applications. The command below is an example to
start `autharr` of the container image `volkerraschek/tarr` via docker. `autharr` is watching for changes of the API
token. Any change will be written to the standard output.

> [!NOTE]
> Adapt the volume mount, if you want to write the token to file on the host system.

```bash
$ docker run \
  --rm \
  --volume /etc/bazarr:/etc/bazarr:ro \
    docker.io/volkerraschek/tarr:latest \
      autharr --watch /etc/bazarr/config.yaml
```

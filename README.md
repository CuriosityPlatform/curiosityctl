# Curosityctl

A tool to manage, init curiosity environment(local, dev, prod)

## About 
 
This tool interacts with [Platform](https://github.com/CuriosityPlatform/platform) to manage curiosity platform environment

## Supported features

All supported commands you can find with curiosityctl --help

Commands shortlist

 - `up` - start environment - creates containers, build images, migrate main database
 - `down` - stop environment
 - `restart` - restarts environment - launches down and up


## Dependencies

### Project

 - `golang` - 1.16
 - `docker`
 - `docker compose` - a docker plugin

### Build dependencies

 - `docker` - v20.10.7


## Installing and developing

### Build project

It will run build, linter, unit-tests

```shell
make
```


### Only build

```shell
make build
```

### Run linter

```shell
make check
```

### Run tests

```shell
make test
```

### Build cache

Docker builder collects **cache**, to clear that cache use
```shell
make cache-clear
```
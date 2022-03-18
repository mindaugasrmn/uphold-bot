# uphold-bot
Uphold bot


## Go Commands


###  Build
```sh
make build
```
builds binary file

### Run

```sh
make run
```
runs main.go 

### Test
```sh
make tests
```
runs tests 

## Docker
###  Build
```sh
make docker-build
```
builds docker image

### Start
```sh
make docker-run
```
runs docker image

### Enter
```sh
make docker-list
```
lists running docker images
```sh
make docker-enter id="<docker_id>"
```
enters docker image with id
```sh
./main
```
run app from docker container

### Stop
```sh
make docker-list
```
lists running docker images

```sh
make docker-stop id="<docker_id>"
```
stops docker container by id


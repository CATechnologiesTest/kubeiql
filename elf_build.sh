# Build our executable in a docker container to make sure we get a linux/ELF
# binary.  Useful for development on a mac.  Not sure about windows...
APP_NAME=kubeiql.elf
docker run --rm -e "GOPATH=/usr" -e "CGO_ENABLED=0" -v "$PWD":/usr/src/${APP_NAME} -w /usr/src/${APP_NAME} golang:1.10 go build

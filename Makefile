# This is how we want to name the binary output
TARGET=dayan-community-srv-proposal

# These are the values we want to pass for Version and BuildTime
GITTAG=`git describe --tags`
BUILD_TIME=`date +%FT%T%z`

# Setup the -ldflags option for go build here, interpolate the variable values
# LDFLAGS=-ldflags "-X main.GitTag=${GITTAG} -X main.BuildTime=${BUILD_TIME}"
LDFLAGS=-ldflags "-X main.BuildTime=${BUILD_TIME}"

.PHONY:all clean release docker
all:clean release

clean:
	rm -f ${TARGET}

release:
	rm -f ${TARGET} && CGO_ENABLED=0 go build ${LDFLAGS} -o ${TARGET} main.go

docker:
	docker build --build-arg CONFIG_HOST=$CONFIG_HOST TARGET=$TARGET -t $TARGET:latest






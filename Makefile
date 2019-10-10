GITCOMMIT=$(shell git rev-parse HEAD)
BUILDTIME=$(shell date)

REMOTE=github.com
USER=turbonomic
PROJECT=turbotower
PROJECTPATH=${REMOTE}/${USER}/${PROJECT}
BINARY=turbotower
OUTPUT_DIR=build

build: clean
	go build -x -ldflags "-X '${PROJECTPATH}/version.GitCommit=${GITCOMMIT}' -X '${PROJECTPATH}/version.BuildTime=${BUILDTIME}'" -o ${OUTPUT_DIR}/${BINARY} 

.PHONY: clean
clean:
	rm -rf ${OUTPUT_DIR}

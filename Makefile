GITCOMMIT=$(shell git rev-parse HEAD)
BUILDTIME=$(shell date -R)

REMOTE=github.com
USER=turbonomic
PROJECT=turbotower
PROJECTPATH=${REMOTE}/${USER}/${PROJECT}
BINARY=turbotower

build: clean
	go build -ldflags "-X '${PROJECTPATH}/version.GitCommit=${GITCOMMIT}' -X '${PROJECTPATH}/version.BuildTime=${BUILDTIME}'" -o ${BINARY} 

.PHONY: clean
clean:
	@rm -rf ${BINARY}

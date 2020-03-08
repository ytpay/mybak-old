BUILD_VERSION    := $(shell cat version)
BUILD_DATE       := $(shell date "+%F %T")
COMMIT_SHA1      := $(shell git rev-parse HEAD)
MYSQL_USER		 := $(MYSQL_USER)
MYSQL_PASSWORD	 := $(MYSQL_PASSWORD)
MYBAK_SECRET     := $(MYBAK_SECRET)
MYBAK_COMMENT    := $(MYBAK_COMMENT)

all: clean
	gox -osarch="darwin/amd64 linux/amd64" \
		-output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" \
		-ldflags	"-X 'main.Version=${BUILD_VERSION}' \
					-X 'main.BuildDate=${BUILD_DATE}' \
					-X 'main.CommitID=${COMMIT_SHA1}'\
					-X 'main.User=${MYSQL_USER}'\
					-X 'main.Password=${MYSQL_PASSWORD}'\
					-X 'main.Secret=${MYBAK_SECRET}' \
					-X 'main.Comment=${MYBAK_COMMENT}'"

clean:
	rm -rf dist

install:
	go install -ldflags		"-X 'main.Version=${BUILD_VERSION}' \
               				-X 'main.BuildDate=${BUILD_DATE}' \
               				-X 'main.CommitID=${COMMIT_SHA1}'\
               				-X 'main.User=${MYSQL_USER}'\
               				-X 'main.Password=${MYSQL_PASSWORD}'\
               				-X 'main.Secret=${MYBAK_SECRET}' \
               				-X 'main.Comment=${MYBAK_COMMENT}'"

.PHONY: all release clean install

.EXPORT_ALL_VARIABLES:

GO111MODULE = on
GOPROXY = https://goproxy.io
GOSUMDB = sum.golang.google.cn

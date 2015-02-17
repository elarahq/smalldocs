PID      = /tmp/smalldocs.pid
GO_FILES = $(wildcard *.go)

serve:
	@make restart
	@fswatch -o . | xargs -n1 -I{}  make restart || make kill

kill:
	@kill `cat $(PID)` || true

restart:
	@make kill
	@go run $(GO_FILES) & echo $$! > $(PID)

install:
	@bower install
	@go get

build:
	@jsx public/js/ public/build/

start:
	@go run $(GO_FILES)

.PHONY: serve restart kill install build start # let's go to reserve rules names

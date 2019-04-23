.PHONY: all
.PHONY: write_api
.PHONY: read_api
.PHONY: consumer

all:
	make consumer
	make write_api
	make read_api

write_api:
	cd write_api && go build && docker build .

read_api:
	cd read_api && go build && docker build .

consumer:
	cd consumer/checksum && go build && docker build .
	cd consumer/delete_event && go build && docker build .
	cd consumer/post_event && go build && docker build .

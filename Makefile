project_path := $(shell git rev-parse --show-toplevel)

dev:	$(project_path)/main.go
	cd $(project_path) && \
	buffalo dev

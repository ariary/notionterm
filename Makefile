before.build:
	go mod tidy -compat=1.17 && go mod download

build.notionterm:
	@echo "build in ${PWD}";go build -o notionterm notionterm.go
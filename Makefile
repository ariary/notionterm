before.build:
	go mod tidy && go mod download && go mod vendor

build.notionterm:
	@echo "build in ${PWD}";go build -o notionterm notionterm.go
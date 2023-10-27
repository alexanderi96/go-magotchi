APPNAME := go-magotchi
REPO := github.com/alexanderi96/go-magotchi

LDFLAGS_BEEPBERRY := -ldflags="-X 'main.scrWdt=400' -X 'main.scrHgt=240'"

.PHONY: clean
clean:
	rm go-magotchi

.PHONY: build
build:
	go build -o go-magotchi

.PHONY: build-beepberry
build-beepberry:
	go build -o go-magotchi ${LDFLAGS_BEEPBERRY}

.PHONY: run
run:
	go run .

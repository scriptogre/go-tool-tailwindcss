build:
	go build -o tailwindcss .

install: build
	sudo cp tailwindcss $$(go env GOROOT)/pkg/tool/$$(go env GOOS)_$$(go env GOARCH)/

clean:
	rm -f tailwindcss

.PHONY: build install clean
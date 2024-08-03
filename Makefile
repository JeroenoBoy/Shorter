dev:
	make -j5 dev/tailwindcss dev/templ dev/air

dev/tailwindcss:
	tailwindcss -i view/styles/index.css -o static/index.css --watch

dev/templ:
	templ generate --watch

dev/air:
	air

build:
	tailwindcss -i view/styles/index.css -o static/index.css --minify
	templ generate
	rm -r out | true
	mkdir out
	cp -r static out/static
	go build -o out/shorter main.go
	

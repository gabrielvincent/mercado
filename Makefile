build_templ:
	~/go/bin/templ generate

build_tailwindcss:
	npx tailwindcss -i base.css -o public/css/tailwind.css --minify

build: build_templ build_tailwindcss
	go build -o mercado main.go


watch_tailwindcss:
	npx tailwindcss -i base.css -o public/css/tailwind.css --watch

watch_go:
	wgo -file=.go -file=.templ -xfile=_templ.go templ generate :: go run main.go

dev: export STAGE=dev
dev:
	npx concurrently --kill-others --raw "make watch_tailwindcss" "make watch_go"
	go run main.go

debug:
	dlv attach $$(pgrep mercado)

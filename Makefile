build_templ:
	~/go/bin/templ generate

build_tailwindcss:
	npx tailwindcss -i base.css -o public/css/tailwind.css --minify

build: build_templ build_tailwindcss
	go build -o mercado main.go

watch_templ:
	templ generate --watch

watch_tailwindcss:
	npx tailwindcss -i base.css -o public/css/tailwind.css --watch

watch_go:
	~/go/bin/air

dev: export STAGE=dev
dev:
	npx concurrently --kill-others --raw "make watch_tailwindcss" "make watch_templ" "make watch_go"
	go run main.go

debug:
	dlv attach $$(pgrep mercado)

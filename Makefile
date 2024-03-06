build:
	npx tailwindcss -i base.css -o public/css/tailwind.css --minify
	go build main.go

dev:
	~/go/bin/air
	go run main.go

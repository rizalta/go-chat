templ:
	@templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v

server:
	@air \
	--build.cmd "go build -o ./tmp/go-chat ./cmd/api/main.go" --build.bin "./tmp/go-chat" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go,templ" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

tailwind:
	@./tailwindcss -i cmd/web/assets/css/input.css -o cmd/web/assets/css/output.css --watch

sync_assets:
	@air \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "cmd/web/assets" \
	--build.include_ext "js,css"

db-up:
	docker-compose up -d

db-down:
	docker-compose down

watch: 
	make -j3 tailwind templ sync_assets

.PHONY: watch tailwind server templ db-up db-down

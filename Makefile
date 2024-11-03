live/templ:
	@templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v

live/server:
	@air \
	--build.cmd "go build -o bin/go-chat cmd/api/main.go" --build.bin "bin/go-chat" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go,templ" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

live/tailwind:
	@./tailwindcss -i cmd/web/assets/css/input.css -o cmd/web/assets/css/output.css --watch --minify

live: 
	make -j4 live/tailwind live/templ live/server

.PHONY: live live/tailwind live/server live/templ

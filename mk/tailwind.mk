TAILWIND_BIN=./tailwindcss

.PHONY: setup-tailwind build vercel-build

setup-tailwind:
	@echo "Downloading tailwind css cli..."
	@curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -o $(TAILWIND_BIN)
	@chmod +x $(TAILWIND_BIN)

build: setup-tailwind
	rm -rf dist && \
	go run ./cmd/ssg && \
	$(TAILWIND_BIN) -i internal/templates/input.css -o dist/styles.css --minify; \
	rm $(TAILWIND_BIN);

vercel-build: setup-go setup-tailwind
	@export PATH=$(PWD)/$(GO_DIR)/go/bin:$$PATH; \
	go run ./cmd/ssg && \
	if [ -f $(TAILWIND_BIN) ]; then \
		$(TAILWIND_BIN) -i internal/templates/input.css -o dist/styles.css --minify; \
		rm $(TAILWIND_BIN); \
	fi

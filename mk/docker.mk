DEV_IMAGE=mehub-dev
DEV_CONTAINER=mehub-dev

.PHONY: dev-build dev-run dev-clean

## dev-build: Build the local development container image.
dev-build:
	podman build -f Dockerfile.dev -t $(DEV_IMAGE) .

## dev-run: Run an interactive shell inside the dev container with the repo mounted.
dev-run:
	podman run --rm -it \
		-v "$(PWD)":/workspace:Z \
		-w /workspace \
		-p 8080:8080 \
		--name $(DEV_CONTAINER) \
		$(DEV_IMAGE)

## dev-clean: Remove the local development container image.
dev-clean:
	podman rmi --force $(DEV_IMAGE)
	@echo "Image '$(DEV_IMAGE)' removed."

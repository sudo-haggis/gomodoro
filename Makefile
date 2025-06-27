# Add this to your existing Makefile, right after the help target:

## build: Build the application
build:
	go build -o $(BINARY_NAME) .

# Also make sure your install target looks like this:
## install: Install GoModoro system-wide (requires sudo)
install: build
	@if [ ! -f "./install.sh" ]; then \
		echo "❌ install.sh not found. Please create it first."; \
		exit 1; \
	fi
	@if [ "$(shell id -u)" != "0" ]; then \
		echo "❌ Installation requires sudo. Run: sudo make install"; \
		exit 1; \
	fi
	./install.sh

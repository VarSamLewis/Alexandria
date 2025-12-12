.PHONY: build install clean uninstall test

# Binary name
BINARY_NAME=alexandria

# Installation directory
INSTALL_DIR=$(HOME)/.local/bin

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME)
	@echo "Build complete: $(BINARY_NAME)"

# Install the binary to ~/.local/bin
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	install -Dm755 $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installation complete!"
	@echo "Make sure $(INSTALL_DIR) is in your PATH"
	@echo "Add this to your ~/.bashrc or ~/.zshrc if needed:"
	@echo "  export PATH=\"\$$HOME/.local/bin:\$$PATH\""

# Remove the binary from installation directory
uninstall:
	@echo "Removing $(BINARY_NAME) from $(INSTALL_DIR)..."
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Uninstall complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Build and run
run: build
	./$(BINARY_NAME)

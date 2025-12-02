# Odek UI Toolkit Makefile

ODIN := odin
BUILD_DIR := build

DEBUG_FLAGS := -debug
RELEASE_FLAGS := -o:speed -disable-assert

EXAMPLES := catalog todo filebrowser

.PHONY: all release test clean help $(EXAMPLES) $(addsuffix -release,$(EXAMPLES))

# Default: build all examples in debug mode
all: $(EXAMPLES)

# Build all examples in release mode
release: $(addsuffix -release,$(EXAMPLES))

# Debug builds
catalog: | $(BUILD_DIR)
	$(ODIN) build examples/catalog -out:$(BUILD_DIR)/catalog $(DEBUG_FLAGS)

todo: | $(BUILD_DIR)
	$(ODIN) build examples/todo -out:$(BUILD_DIR)/todo $(DEBUG_FLAGS)

filebrowser: | $(BUILD_DIR)
	$(ODIN) build examples/filebrowser -out:$(BUILD_DIR)/filebrowser $(DEBUG_FLAGS)

# Release builds
catalog-release: | $(BUILD_DIR)
	$(ODIN) build examples/catalog -out:$(BUILD_DIR)/catalog-release $(RELEASE_FLAGS)

todo-release: | $(BUILD_DIR)
	$(ODIN) build examples/todo -out:$(BUILD_DIR)/todo-release $(RELEASE_FLAGS)

filebrowser-release: | $(BUILD_DIR)
	$(ODIN) build examples/filebrowser -out:$(BUILD_DIR)/filebrowser-release $(RELEASE_FLAGS)

# Run tests
test:
	$(ODIN) test tests

# Create build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Help
help:
	@echo "Odek UI Toolkit Build Targets:"
	@echo ""
	@echo "  make              Build all examples (debug)"
	@echo "  make release      Build all examples (release)"
	@echo ""
	@echo "  make catalog      Build catalog example (debug)"
	@echo "  make todo         Build todo example (debug)"
	@echo "  make filebrowser  Build filebrowser example (debug)"
	@echo ""
	@echo "  make catalog-release      Build catalog (release)"
	@echo "  make todo-release         Build todo (release)"
	@echo "  make filebrowser-release  Build filebrowser (release)"
	@echo ""
	@echo "  make test         Run test suite"
	@echo "  make clean        Remove build artifacts"
	@echo "  make help         Show this help"

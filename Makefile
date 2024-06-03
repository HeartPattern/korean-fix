PLATFORMS := \
	darwin/amd64 \
	darwin/arm64 \
	linux/amd64 \
	linux/arm \
	linux/arm64 \
	windows/amd64 \
	windows/386

MAIN_FILE := ./cmd/main.go

BUILD_DIR := build

os_arch = $(subst /, ,$1)
os = $(word 1,$(call os_arch,$1))
arch = $(word 2,$(call os_arch,$1))

all: $(PLATFORMS)

$(PLATFORMS):
	@mkdir -p $(BUILD_DIR)
	GOOS=$(call os,$@) GOARCH=$(call arch,$@) go build -o $(BUILD_DIR)/korean-fix-$(call os,$@)-$(call arch,$@) $(MAIN_FILE)
	@echo "Built $(BINARY_NAME) for $(call os,$@)/$(call arch,$@)"

clean:
	rm -rf $(BUILD_DIR)

.PHONY: all clean $(PLATFORMS)

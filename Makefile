# Build
build:
	@echo "Build proto..."
	@protoc --go_out=. ./internal/core/proto/health/*.proto
	@protoc --go_out=. ./internal/core/proto/payment/*.proto
	@protoc --go_out=. ./internal/core/proto/pod/*.proto

.PHONY: build
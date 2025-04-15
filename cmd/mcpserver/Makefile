APP_NAME := mcpserver
IMAGE_NAME := napat/mcpserver-demo
IMAGE_TAG := 1.0
# MODEL := google:gemini-2.5-pro-exp-03-25
MODEL := google:gemini-2.0-flash-001

.PHONY: mcp-docker-run
mcp-docker-run:
	@echo "Running MCP Server locally..."
	@echo "Loading API key from environment..."
	bash -c '. ./export-key.sh && mcphost --google-api-key "$$GOOGLE_API_KEY" --model $(MODEL) --config ./napat-mcp-server.json'
	
.PHONY: mcp-docker-build
mcp-docker-build:
	@echo "Building MCP Server Docker image..."
	cd ../.. && docker build --no-cache -t $(IMAGE_NAME):$(IMAGE_TAG) -f cmd/mcpserver/Dockerfile .

.PHONY: mcp-docker-push
mcp-docker-push:
	@echo "Pushing MCP Server Docker image to hub.docker.com/repositories/napat..."
	docker push $(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: mcp-help
mcp-help:
	@echo "MCP Server Makefile"
	@echo "--------------------"
	@echo "Available commands:"
	@echo "  mcp-docker-build    - Build the MCP Server Docker image"
	@echo "  mcp-docker-push     - Push Docker image to Docker Hub"
	@echo "  mcp-help     - Show this help message" 
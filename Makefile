CONFIG_FILE = config.json
ENV_FILE = .env
OS_NAME := $(shell uname -s)

# Cek dan install jq jika belum ada
check-jq:
	@echo "🔍 Checking for jq..."
	@if ! command -v jq >/dev/null 2>&1; then \
		echo "⚠️  jq not found. Installing..."; \
		if [ "$(OS_NAME)" = "Linux" ]; then \
			if command -v apt >/dev/null 2>&1; then \
				sudo apt update && sudo apt install -y jq; \
			elif command -v dnf >/dev/null 2>&1; then \
				sudo dnf install -y jq; \
			else \
				echo "❌ Unsupported Linux package manager. Please install jq manually."; exit 1; \
			fi; \
		elif [ "$(OS_NAME)" = "Darwin" ]; then \
			if command -v brew >/dev/null 2>&1; then \
				brew install jq; \
			else \
				echo "❌ Homebrew not found. Please install jq manually: https://brew.sh"; exit 1; \
			fi; \
		else \
			echo "❌ Unsupported OS ($(OS_NAME)). Please install jq manually."; exit 1; \
		fi; \
	else \
		echo "✅ jq found."; \
	fi


.PHONY: env up down restart logs clean help check-jq

env: check-jq
	@echo "🔧 Generating $(ENV_FILE) from $(CONFIG_FILE)..."
	@if [ ! -f $(CONFIG_FILE) ]; then \
		echo "❌ $(CONFIG_FILE) not found!"; exit 1; \
	fi
	@rm -f $(ENV_FILE)
	@jq -r '.postgres | to_entries[] | "POSTGRES_\(.key|ascii_upcase)=\(.value)"' $(CONFIG_FILE) >> $(ENV_FILE)
	@jq -r '.redis | to_entries[] | "REDIS_\(.key|ascii_upcase)=\(.value)"' $(CONFIG_FILE) >> $(ENV_FILE)
	@jq -r '.minio | to_entries[] | "MINIO_\(.key|ascii_upcase)=\(.value)"' $(CONFIG_FILE) >> $(ENV_FILE)
	@jq -r '.server | to_entries[] | "SERVER_\(.key|ascii_upcase)=\(.value)"' $(CONFIG_FILE) >> $(ENV_FILE)
	@echo "✅ $(ENV_FILE) generated successfully!"

up: env
	@echo "🚀 Starting containers..."
	@docker compose up -d
	@echo "✅ All services are running."

# Stop semua container
down:
	@echo "🛑 Stopping containers..."
	@docker compose down
	@echo "✅ All containers stopped."

restart: down up

pull:
	@echo "⬇️  Pulling latest images..."
	@docker compose pull
	@echo "✅ Images updated."

build:
	@echo "🏗️  Building images..."
	@docker compose up -d --build
	@echo "✅ Build complete."


logs:
	@docker compose logs -f

clean:
	@echo "🧹 Cleaning up containers, networks, and volumes..."
	@docker compose down -v --remove-orphans
	@rm -f $(ENV_FILE)
	@echo "✅ Clean complete."

help:
	@echo ""
	@echo "🧩 Available commands:"
	@echo "  make env       → Generate .env from config.json"
	@echo "  make up        → Generate env & start containers"
	@echo "  make down      → Stop containers"
	@echo "  make restart   → Restart all containers"
	@echo "  make logs      → Tail logs from all containers"
	@echo "  make clean     → Remove containers, volumes & .env"
	@echo ""

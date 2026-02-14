.PHONY: dev stop

dev:
	docker compose up -d
	@echo "Waiting for Postgres..."
	@until docker compose exec db pg_isready -U postgres > /dev/null 2>&1; do sleep 0.5; done
	@echo "Postgres is ready."
	mix ecto.create 2>/dev/null; true
	mix phx.server

stop:
	docker compose down

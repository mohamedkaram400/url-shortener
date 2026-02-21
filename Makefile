migration-up:
	migrate -database "$(DB_URL)" -path ./migrations up

migration-down:
	migrate -database "$(DB_URL)" -path ./migrations down 1

migration-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migration-create name=create_users"; \
		exit 1; \
	fi
	migrate create -ext sql -dir ./migrations -format unix $(name)
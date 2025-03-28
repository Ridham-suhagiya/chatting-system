DB_URL ?= $(shell printenv DB_URL)

migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" -verbose down

DB_URL="postgresql://localhost:5432/chatting-system?sslmode=disable"

migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" -verbose down

up:
	docker compose up -d

db:
	docker compose exec db psql -U postgres -d docker

# $USAGE ---------
# make migration.new name=create_users
migration.new:
	docker compose exec api goose -dir migrations create $(name).sql sql

migration.up:
	docker compose exec api goose -dir migrations postgres "user=postgres dbname=docker host=db password=password sslmode=disable" up

migration.down:
	docker compose exec api goose -dir migrations postgres "user=postgres dbname=docker host=db password=password sslmode=disable" down

migration.status:
	docker compose exec api goose -dir migrations postgres "user=postgres dbname=docker host=db password=password sslmode=disable" status

dump.schema:
	docker compose exec db pg_dump -U postgres -d docker --schema-only > schema.sql

sqlc.generate:
	docker compose exec api sqlc generate

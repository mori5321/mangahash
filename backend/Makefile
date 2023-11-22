up:
	docker compose up -d

db:
	docker compose exec -e PGOPTIONS="--search_path=app"  db psql -U postgres -d docker;

# $USAGE ---------
# make migration.new name=create_users
migration.new:
	docker compose exec api goose -dir migrations create $(name).sql sql

migration.up:
	docker compose exec api goose -dir migrations postgres "user=postgres dbname=docker host=db password=password sslmode=disable" up
	$(MAKE) dump.schema

migration.down:
	docker compose exec api goose -dir migrations postgres "user=postgres dbname=docker host=db password=password sslmode=disable" down
	$(MAKE) dump.schema

migration.status:
	docker compose exec api goose -dir migrations postgres "user=postgres dbname=docker host=db password=password sslmode=disable" status

test:
	$(MAKE) test.db.restore
	docker compose run api go test -v ./...

test.db:
	docker compose exec -e PGOPTIONS="--search_path=app" testdb psql -U test -d test

test.migration.reset:
	docker compose exec api goose -dir migrations postgres "user=test dbname=test host=testdb password=password sslmode=disable" reset


test.migration.up:
	docker compose exec api goose -dir migrations postgres "user=test dbname=test host=testdb password=password sslmode=disable" up

test.migration.status:
	docker compose exec api goose -dir migrations postgres "user=test dbname=test host=testdb password=password sslmode=disable" status

test.db.restore:
	$(MAKE) test.migration.reset
	$(MAKE) test.migration.up 

schema.dump:
	docker compose exec db pg_dump -U postgres -d docker --schema app --schema-only > sqlc/schema.sql

sqlc.gen:
	docker compose exec api sqlc generate

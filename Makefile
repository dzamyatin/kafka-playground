.PHONY: schema-log,schema-sh,up,down

include .env

schema-sh:
	docker exec -ti $(PROJECT_NAME)-schema-registry sh
schema-log:
	docker logs -f $(PROJECT_NAME)-schema-registry
up:
	docker compose -f ./docker-compose.yml up --force-recreate
down:
	docker compose -f ./docker-compose.yml down
.PHONY: schema-log,schema-sh,up,down

include .env

schema-sh:
	docker exec -ti $(PROJECT_NAME)-schema-registry sh
schema-log:
	docker logs -f $(PROJECT_NAME)-schema-registry
registry-ui-log:
	docker logs -f $(PROJECT_NAME)-schema-registry-ui
registry-ui-sh:
	docker exec -ti $(PROJECT_NAME)-schema-registry-ui sh
up:
	docker compose -f ./docker-compose.yml up --force-recreate -d
down:
	docker compose -f ./docker-compose.yml down
restart:
	docker compose -f ./docker-compose.yml restart
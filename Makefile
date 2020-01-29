.PHONY: deploy.prd
deploy.prd:
	docker-compose -f docker-compose.prd.yaml pull
	docker-compose -f docker-compose.prd.yaml up -d

.PHONY: setup.server.dev
setup.server.dev:
	docker-compose up -d --build

.PHONY: teardown.server.dev
teardown.server.dev:
	docker-compose down


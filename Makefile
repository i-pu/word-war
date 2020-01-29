.PHONY: deploy.prd
deploy.prd:
	docker-compose -f docker-compose.prd.yaml pull
	docker-compose -f docker-compose.prd.yaml up -d

.PHONY: deploy.dev
deploy.dev:
	docker-compose up -d --build

.PHONY: deploy.prd
deploy.prd:
	docker-compose -f docker-compose.prd.yaml pull
	docker-compose -f docker-compose.prd.yaml up -d

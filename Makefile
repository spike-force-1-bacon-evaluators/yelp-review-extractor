.DEFAULT_GOAL := help

PROJECT_NAME := Yelp - Review Extractor
PROJECT_ALIAS := bacon
IMAGE_NAME := yelp-review-extractor

.PHONY: help build test clean

help:
	@echo "------------------------------------------------------------------------"
	@echo "${PROJECT_NAME}"
	@echo "------------------------------------------------------------------------"
	@grep -E '^[a-zA-Z_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: clean ## Build unit test container
	@docker build -t ${PROJECT_ALIAS}/${IMAGE_NAME} -f resources/Dockerfile .

test: build ## Run unit tests
	@docker run --rm ${PROJECT_ALIAS}/${IMAGE_NAME}

clean: ## Remove images and containers
	@./resources/scripts/rm-image.sh ${PROJECT_ALIAS}/${IMAGE_NAME}
	@./resources/scripts/rm-container.sh ${IMAGE_NAME}

VERSION := v1.0.0
IMAGE_TAG := "asia.gcr.io/presstotalk/api:${VERSION}"


.PHONY: build
build:
	@docker build -t ${IMAGE_TAG} .
	@docker push ${IMAGE_TAG}


.PHONY: deploy
deploy:
	@gcloud run deploy api --image ${IMAGE_TAG}

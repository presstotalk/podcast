VERSION := v1.0.0
IMAGE_TAG := "asia.gcr.io/presstotalk/api:${VERSION}"


.PHONY: build
build:
	@docker buildx build --platform linux/amd64 -t ${IMAGE_TAG} .
	@docker push ${IMAGE_TAG}


.PHONY: deploy
deploy:
	@gcloud run deploy api --region=asia-east1 --image=${IMAGE_TAG}

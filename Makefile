VERSION := v1.0.0
IMAGE_TAG := "asia.gcr.io/presstotalk/api:${VERSION}"


.PHONY: run
run:
	@go run cmd/server/main.go


.PHONY: build
build:
	@docker buildx build --platform linux/amd64 -t ${IMAGE_TAG} .
	@docker push ${IMAGE_TAG}


.PHONY: deploy
deploy:
	@gcloud run deploy api \
		--no-use-http2 \
		--region=asia-east1 \
		--env-vars-file=.env \
		--image=${IMAGE_TAG}

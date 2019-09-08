docker:
	drone exec
	docker build -t jsussman/drone-validator:${VERSION} .
	docker push jsussman/drone-validator:${VERSION}

install:
	@go get -u github.com/golang/lint/golint
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o person-api .
	docker build -t person-api .
	docker tag person-api ${USERNAME}/person-api:1.0.0
	docker rmi person-api
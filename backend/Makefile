.PHONY: test
test:
	go test ./... -coverprofile cover.out

build:
	CGO_ENABLED=0 go build -a -o ./bin/qe-dashboard-backend ./cmd/qe-dashboard-backend/*

fmt:
	gofmt -l -s -w ./
	goimports -l -w ./

swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	go get github.com/swaggo/swag/gen@latest
	go get github.com/swaggo/swag/cmd/swag@latest
	cd pkg/api && $$(go env GOPATH)/bin/swag init -g server.go

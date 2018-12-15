all: build

deploy:
	scp winona chip@wbez-button.local:winona

build:
	@GOARM=7 GOARCH=arm GOOS=linux go build

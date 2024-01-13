.PHONY: build run

build:
	docker build -t stress-test .

run: build
	clear
	docker run -ti stress-test ./build/main —url=http://google.com —requests=1000 —concurrency=100
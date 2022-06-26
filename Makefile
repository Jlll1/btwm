build: clean
	go build -o btwm main.go

run: build
	echo "exec ./btwm" > xinitrc
	./run.sh

clean:
	go clean
	rm -f xinitrc

.PHONY: build run clean
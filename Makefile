btwm: main.go
	go build -o btwm main.go

run: btwm
	echo "exec btwm" > xinitrc
	./run.sh

clean:
	rm btwm xinitrc

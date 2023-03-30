.PHONY: run

dev: clean run

run:
	go run . generate \
	--name github.com/mahcks/test-project \
	--path /Users/mahcks/Desktop/Stack-Test \
	--mariadb \
	--redis

clean:
	rm -rf /Users/mahcks/Desktop/Stack-Test/*

comma:=,
space:=\ 
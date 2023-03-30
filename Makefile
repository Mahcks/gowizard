.PHONY: run

MODULE := github.com/someone/module
MODULE_PATH := /module/path/somewhere

dev: clean run

run:
	go run . generate \
	--name $(MODULE) \
	--path $(MODULE_PATH) \
	--mariadb \
	--redis

clean:
	rm -rf $(MODULE_PATH)/*

comma:=,
space:=\ 
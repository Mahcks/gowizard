.PHONY: run

MODULE := github.com/someone/module
MODULE_PATH := /module/path/somewhere

dev: clean run

run:
	go run . generate \
	--module $(MODULE) \
	--path $(MODULE_PATH) \
	--adapter mariadb,redis,mongodb

clean:
	rm -rf $(MODULE_PATH)/*

comma:=,
space:=\ 
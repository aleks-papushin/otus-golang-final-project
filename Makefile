.PHONY: linux
linux:
	docker build -t stat_collector .
	docker run --rm stat_collector
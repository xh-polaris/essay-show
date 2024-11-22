.PHONY: start build wire clean

run:
	sh ./output/bootstrap.sh
build:
	sh ./build.sh
build_and_run:
	sh ./build.sh && sh ./output/bootstrap.sh
wire:
	wire ./provider
clean:
	rm -r ./output

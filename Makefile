.PHONEY: example

.DEFAULT_GOAL := example

clean:
	rm -rf vp

vp:
	go build -o vp cmd/vp/main.go

example: clean vp
	./vp test.jpg

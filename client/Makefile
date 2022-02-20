build:
	mkdir bin
	PKG_CONFIG_PATH="/usr/local/opt/python@3.7/lib/pkgconfig" LIBRARY_PATH="/usr/local/opt/python@3.7/lib" CGOENABLED=1 go build -o bin/pydc_client.so -buildmode=c-shared .

clean:
	rm -rf bin/

test:
	go test ./...

pytest:
	python3 tester.py

all:
	make clean
	make build
	make pytest

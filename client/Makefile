build:
	CGOENABLED=1 go build -o pydc_client.so -buildmode=c-shared .

clean:
	rm ~/percy/pydc/client/*.so
	rm ~/percy/pydc/client/*.h

test:
	go test ./...

pytest:
	python3 tester.py

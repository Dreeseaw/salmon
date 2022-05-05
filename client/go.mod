module github.com/Dreeseaw/salmon/client

go 1.17

require (
	github.com/Dreeseaw/salmon/shared/grpc v0.0.0-20220505012613-79dc0ed87b6a
	github.com/google/uuid v1.3.0
	github.com/kelindar/column v0.0.0-20220310063741-ee265e7d894c
	github.com/stretchr/testify v1.7.1
	google.golang.org/grpc v1.46.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/kelindar/bitmap v1.1.5 // indirect
	github.com/kelindar/intmap v1.1.0 // indirect
	github.com/kelindar/iostream v1.3.0 // indirect
	github.com/kelindar/smutex v1.0.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/zeebo/xxh3 v1.0.1 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/Dreeseaw/salmon/shared/grpc => ../shared/grpc

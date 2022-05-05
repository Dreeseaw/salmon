module github.com/Dreeseaw/salmon/examples

go 1.17

require github.com/Dreeseaw/salmon/client v0.0.0-20220424153844-4c66ba1c5c55

require (
	github.com/Dreeseaw/salmon/grpc v0.0.0-20220424153844-4c66ba1c5c55 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/kelindar/bitmap v1.2.1 // indirect
	github.com/kelindar/column v0.0.0-20220310063741-ee265e7d894c // indirect
	github.com/kelindar/intmap v1.1.0 // indirect
	github.com/kelindar/iostream v1.3.0 // indirect
	github.com/kelindar/smutex v1.0.0 // indirect
	github.com/klauspost/compress v1.15.2 // indirect
	github.com/klauspost/cpuid/v2 v2.0.12 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
	golang.org/x/sys v0.0.0-20220429233432-b5fbb4746d32 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220429170224-98d788798c3e // indirect
	google.golang.org/grpc v1.46.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
    github.com/Dreeseaw/salmon/client => ../../client/
)

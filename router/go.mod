module github.com/Dreeseaw/salmon/router

go 1.17

replace (
	github.com/Dreeseaw/salmon/shared/config => ../shared/config
	github.com/Dreeseaw/salmon/shared/grpc => ../shared/grpc
)

require (
	github.com/Dreeseaw/salmon/shared/config v0.0.0-00010101000000-000000000000
	github.com/Dreeseaw/salmon/shared/grpc v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.46.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

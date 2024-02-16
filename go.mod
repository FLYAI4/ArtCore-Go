module github.com/robert-min/ArtCore-Go

go 1.19

replace github.com/robert-min/ArtCore-Go/src/imagetovideo => ./src/imagetovideo

replace github.com/robert-min/ArtCore-Go/src/pb => ./src/pb

require (
	github.com/joho/godotenv v1.5.1
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/stretchr/testify v1.8.4
	gocv.io/x/gocv v0.35.0
	google.golang.org/grpc v1.61.1
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231106174013-bbf56f31fb17 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

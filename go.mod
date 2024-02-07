module github.com/robert-min/ArtCore-Go

go 1.19

replace github.com/robert-min/ArtCore-Go/src/imagetovideo => ./src/imagetovideo

require (
	github.com/joho/godotenv v1.5.1
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/stretchr/testify v1.8.4
	gocv.io/x/gocv v0.35.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

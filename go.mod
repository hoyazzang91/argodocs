module github.com/hoyazzang91/argodocs

go 1.21

require (
	github.com/spf13/cobra v1.8.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace (
	github.com/hoyazzang91/argodocs/logger => ./logger
	github.com/hoyazzang91/argodocs/markdown => ./markdown
	github.com/hoyazzang91/argodocs/mdgen => ./mdgen
	github.com/hoyazzang91/argodocs/workflow => ./workflow
)

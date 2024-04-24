module github.com/christoofar/safexz

// I wrote this code on 1.22.1, so what I am saying here is that this will work on Go going back to 1.12 (from 2019).
// there were a lot of CVEs addressed by go.dev since that time, but if you are stuck on an older version of Go because of
// your OS/broken package manager, you lost support for your arch, whatever reason, this code does not rely on the newest
// Go levels.
go 1.12

require github.com/stretchr/testify v1.9.0

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

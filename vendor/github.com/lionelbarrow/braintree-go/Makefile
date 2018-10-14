SHELL:=/bin/bash

test:
	go test -parallel 15 -tags='unit integration' ./...

unit:
	go test -tags=unit ./...

integration:
	go test -parallel 15 -tags=integration ./...

analysis:
	diff -u <(echo -n) <(go list -f '{{range .TestGoFiles}}{{$$.ImportPath}}/{{.}}{{end}}' ./...) \
		|| (exit_code=$$?; echo -e '\033[31mTest files should be marked with a unit or integration build tag.\033[0m'; exit $$exit_code)
	diff -u <(echo -n) <(gofmt -d .) \
		|| (exit_code=$$?; echo -e '\033[31mRun gofmt to format source files.\033[0m'; exit $$exit_code)
	go vet ./...
	go get github.com/kisielk/errcheck && CGO_ENABLED=0 errcheck ./...

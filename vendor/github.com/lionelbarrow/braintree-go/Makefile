SHELL:=/bin/bash

test:
	go test -parallel 15 ./...

analysis:
	diff -u <(echo -n) <(gofmt -d .)
	go vet ./...
	go get github.com/kisielk/errcheck && CGO_ENABLED=0 errcheck ./...

analysis-is-backward-compatible-with-master:
	go install \
		&& go get github.com/bradleyfalzon/apicompat/cmd/apicompat \
		&& apicompat -before master ./...

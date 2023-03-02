.PHONY: # Makes every target phony.

# make tidy while you work to catch issues.
tidy: _fix_go _check_go _test

# make sure before pushing to remote.
sure: _fix_go _check_go _integration_test

# should only be run by the cicd pipeline.
_cicd: _deps _check_go _integration_test

# Check Go files.
_check_go:
	go vet ./...
	staticcheck ./...

# Fix Go files, these might alter the files.
_fix_go:
	go mod tidy
	go fmt ./...

_test:
	go test -timeout 1m -v ./...

_integration_test:
	go test --tags=integration -timeout 1m -v ./...

# Intended for CICD use, but could but run locally if dependencies are missing.
_deps:
	go install honnef.co/go/tools/cmd/staticcheck@latest

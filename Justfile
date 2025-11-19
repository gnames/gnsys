# Run the tests of the project
test:
    go test -count=1 -p 1 -shuffle=on -coverprofile=coverage.txt -covermode=atomic ./...

# Run linter (requires golangci-lint)
lint:
    golangci-lint run
  
# List available recipes
list:
    @just --list

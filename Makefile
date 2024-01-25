
benchmarks:
	@echo "Running pkg benchmarks..."
	go test -bench=. -timeout 30m ./pkg

testAll:
	@echo "Running all tests..."
	go test -cover ./...

testRace:
	@echo "Running tests races..."
	go test -race ./...


# Requires godoc to be installed - `go install golang.org/x/tools/cmd/godoc@latest`

serveDocs:
	@echo "Opening in browser..."
	open http://localhost:8989/github.com/Pencroff/ccid_go
	@echo "Serving docs..."
	pkgsite -http=:8989
compile:
	# Linux
	GOOS=linux GOARCH=amd64 go build -o ./bin/hit_linux_amd64 ./cmd/hit
	# OS X
	GOOS=darwin GOARCH=amd64 go build -o ./bin/hit_darwin_amd64 ./cmd/hit
	# OS X M1
	GOOS=darwin GOARCH=arm64 go build -o ./bin/hit_darwin_arm64 ./cmd/hit
	# Windows
	GOOS=windows GOARCH=amd64 go build -o ./bin/hit_windows_amd64.exe ./cmd/hit
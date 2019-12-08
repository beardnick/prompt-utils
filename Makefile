build:build_linux build_mac build_win

clean_linux:
	rm -f bin/linux/*

clean_mac:
	rm -f bin/mac/*

clean_win:
	rm -f bin/win/*

build_linux:clean_linux
	GOOS=linux GOARCH=amd64 go build -v -o bin/linux/redis-prompt prompt-utils/redis
	GOOS=linux GOARCH=amd64 go build -v -o bin/linux/sftp-prompt prompt-utils/sftp


build_mac:clean_mac
	GOOS=darwin GOARCH=amd64 go build -v -o bin/mac/redis-prompt prompt-utils/redis
	GOOS=darwin GOARCH=amd64 go build -v -o bin/mac/sftp-prompt prompt-utils/sftp


build_win:clean_win
	GOOS=windows GOARCH=amd64 go build -v -o bin/win/redis-prompt prompt-utils/redis
	GOOS=windows GOARCH=amd64 go build -v -o bin/win/sftp-prompt prompt-utils/sftp

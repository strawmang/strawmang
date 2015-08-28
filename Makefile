commit=`git rev-parse --short HEAD`
flags=-ldflags "-X main.Commit=$(commit)"
all:
	go build $(flags)
osx:
	GOOS=darwin go build $(flags) -o strawmang-osx
windows:
	GOOS=windows go build $(flags)
linux:
	GOOS=linux go build $(flags)
	

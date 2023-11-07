flags="-X 'main.goVersion=$(go version)' -X main.buildTime=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.commit=`git describe --long --dirty --abbrev=14`"
go build -ldflags "$flags"  -o web_hook main.go
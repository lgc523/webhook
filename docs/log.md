## gin 

```go
go get -u github.com/gin-gonic/gin

go get github.com/fsnotify/fsnotify
go get github.com/spf13/viper

go get -u github.com/jinzhu/gorm
go get -u github.com/go-sql-driver/mysql


go get -u kafka "github.com/segmentio/kafka-go"

go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp

go get github.com/stretchr/testify

go get github.com/bytedance/sonic
```

```shell
sonar-scanner \
  -Dsonar.projectKey=gin_webhook \
  -Dsonar.sources=. \
  -Dsonar.host.url=xxx \
  -Dsonar.token=xxx \
  -Dsonar.go.coverage.reportPaths=coverage.out 
```

```shell
``test with coverage``
go install github.com/axw/gocov/gocov@latest
go install github.com/AlekSi/gocov-xml@latest

go test ./... -coverprofile=coverage.out
gocov convert coverage.out | gocov-xml > coverage.xml
go tool cover -html=coverage.out -o coverage.html
```

```shell
go build -ldflags="-s -w" -o
-s removes the symbol table from the executable to reduce binary size.
-w disables DWARF generation, which also reduces binary size.
```

```shell
go install golang.org/x/vuln/cmd/govulncheck@latest
```

```shell
input {
	file{
      		path => "/var/lib/mysql/v2.log"
      		start_position => "end"
		sincedb_path => "/dev/null"
		codec => multiline {
			 pattern => "^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{6}Z"
		         negate => true
		         what => "previous"
    		}
   	}
}
filter {
  grok {
    match => { "message" => "%{TIMESTAMP_ISO8601:timestamp} %{INT:thread_id} %{WORD:query_type}" }
  }
}

output {
  kafka {
    bootstrap_servers => "xxx"
    topic_id => "t-mysql-general-log"
    codec => json
  }
}
```


```shell
go mod why -m dep
```
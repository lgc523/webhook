all: test lint sonar clean upload

test:
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@gocov convert coverage.out | gocov-xml > coverage.xml

lint: 
	@golangci-lint run

sonar:
	@sonar-scanner \
  	-Dsonar.projectKey=gin_webhook \
  	-Dsonar.sources=. \
  	-Dsonar.host.url=https://sonarqube.youpinsanyue.com \
  	-Dsonar.token=sqp_f62125c296dc7814967a852634d0ff0e45a9506e \
  	-Dsonar.go.coverage.reportPaths=coverage.out \
	-Dsonar.go.golangci-lint.reportPaths=report.xml  

clean:
	@rm -f coverage.out coverage.xml coverage.html


#export COPYFILE_DISABLE=1
upload:
	@tar czf webhook.tar.gz -C .. webhook
	@scp -i ~/.ssh/lgc -P 5235  webhook.tar.gz ubuntu@v1:
	@rm -rf webhook.tar.gz
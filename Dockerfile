FROM golang:1.21-alpine AS builder

WORKDIR  /app

COPY ./ .

ENV GOPROXY https://goproxy.cn

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o webhook main.go

FROM golang:1.21-alpine

ENV SONAR_SERVER https://sonarqube.youpinsanyue.com
ENV SONAR_USERNAME admin
ENV SONAR_PASSWORD UC6pfUvxSTfzD664

WORKDIR /root/
COPY --from=builder /app/webhook .

RUN mkdir conf

COPY --from=builder /app/conf/server.yml conf/

EXPOSE 9031

CMD ["./webhook"]

#RUN export version=$(go version) && \
#    export buildTime=$(date -u '+%Y-%m-%d_%I:%M:%S%p') && \
#    export commit=$(git describe --long --dirty --abbrev=14) && \
#    flags "-X main.goVersion=$version -X main.buildTime=$buildTime -X main.commit=$commit" && \
#    RUN /bin/sh -c 'CGO_ENABLED=0 GOOS=linux time go build -p 4 -ldflags "$flags" -a -o webhook main.go'
#    RUN #CGO_ENABLED=0 GOOS=linux go build -ldflags "$flags" -a -o webhook main.go
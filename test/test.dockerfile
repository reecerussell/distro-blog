FROM golang

COPY . /go/src/github.com/reecerussell/distro-blog

RUN go get github.com/google/uuid
RUN go get github.com/go-sql-driver/mysql
RUN go get golang.org/x/crypto/pbkdf2
RUN go get github.com/aws/aws-lambda-go/events
RUN go get github.com/aws/aws-lambda-go/lambda

WORKDIR /go/src/github.com/reecerussell/distro-blog

ENV AWS_REGION=eu-west-2
ENV AWS_ACCESS_KEY_ID=<key>
ENV AWS_SECRET_ACCESS_KEY=<secret>
ENV JWT_KEY_ID=alias/distro-jwt
ENV CONFIG_BUCKET_NAME=distro-config-store
ENV AUTH_CONFIG_BUCKET_KEY=authorizer-config.yml

CMD sleep 10; go test -v ./... -race -coverprofile=coverage.txt -covermode=atomic

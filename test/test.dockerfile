FROM golang

COPY . /go/src/github.com/reecerussell/distro-blog

RUN go get github.com/google/uuid
RUN go get github.com/go-sql-driver/mysql
RUN go get golang.org/x/crypto/pbkdf2
RUN go get github.com/aws/aws-lambda-go/events
RUN go get github.com/aws/aws-lambda-go/lambda
RUN go get github.com/aws/aws-sdk-go/aws
RUN go get github.com/aws/aws-sdk-go/aws/session
RUN go get github.com/aws/aws-sdk-go/service/kms

WORKDIR /go/src/github.com/reecerussell/distro-blog

ENV AWS_REGION=eu-west-2
ENV AWS_ACCESS_KEY_ID=<aws access id>
ENV AWS_SECRET_ACCESS_KEY=<aws access key>
ENV JWT_KEY_ID=alias/distro-jwt

CMD sleep 15; go test -v ./... -race -coverprofile=coverage.txt -covermode=atomic
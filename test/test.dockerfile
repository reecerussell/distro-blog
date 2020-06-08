FROM golang

COPY . /go/src/github.com/reecerussell/distro-blog

RUN go get github.com/google/uuid
RUN go get github.com/go-sql-driver/mysql
RUN go get golang.org/x/crypto/pbkdf2
RUN go get github.com/aws/aws-lambda-go/events
RUN go get github.com/aws/aws-lambda-go/lambda

WORKDIR /go/src/github.com/reecerussell/distro-blog

CMD sleep 15; go test -v ./... -race -coverprofile=coverage.txt -covermode=atomic
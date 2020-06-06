FROM golang

COPY . /go/src/github.com/reecerussell/distro-blog

RUN go get github.com/google/uuid
RUN go get github.com/go-sql-driver/mysql
RUN go get golang.org/x/crypto/pbkdf2

WORKDIR /go/src/github.com/reecerussell/distro-blog

CMD sleep 20; go test -v ./...
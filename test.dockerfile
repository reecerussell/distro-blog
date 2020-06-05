FROM golang

COPY . /go/src/github.com/reecerussell/distro-blog

RUN go get github.com/google/uuid
RUN go get github.com/go-sql-driver/mysql
RUN go get golang.org/x/crypto/pbkdf2

ENV CONN_STRING=root:password@tcp(db)/distro-blog-test?parseTime=true

WORKDIR /go/src/github.com/reecerussell/distro-blog

CMD sleep 20; go test -v ./...
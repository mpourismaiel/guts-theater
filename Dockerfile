FROM golang:1.17

WORKDIR /src/github.com/mpourismaiel/guts-theater

COPY go.sum /src/github.com/mpourismaiel/guts-theater/
COPY go.mod /src/github.com/mpourismaiel/guts-theater/

RUN echo "Downloading the Go dependencies..." && \
  go mod download && \
  echo "Finished downloading. Building the project..."

COPY . /src/github.com/mpourismaiel/guts-theater/

RUN go build -o guts-theater /src/github.com/mpourismaiel/guts-theater/main.go

EXPOSE 4000
CMD ["/src/github.com/mpourismaiel/guts-theater/guts-theater"]

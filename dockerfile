FROM golang:1.21.6
WORKDIR /banner-service
COPY go.mod go.sum ./
RUN go mod download 
COPY *.go ./
RUN GOOS=linux GOARCH=amd64 go build -o /service
CMD ["service"]
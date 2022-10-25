FROM golang:latest

COPY ./ ./
ENV GOPATH=/
RUN go build cmd/weatherApp.go
CMD ["./weatherApp"]
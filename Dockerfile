FROM golang:1.13 as builder
WORKDIR /code
ADD go.mod go.sum /code/
RUN go mod download
ADD . .
RUN go build -o /topz main.go

FROM gcr.io/distroless/base
ENV VERSION "v0.0.1"
EXPOSE 1080
WORKDIR /
COPY --from=builder /topz /usr/bin/topz
ENTRYPOINT ["/usr/bin/topz"]
FROM golang:1.23.0 AS api

WORKDIR /compiler

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o application ./cmd/app/main.go

FROM scratch AS prod

WORKDIR /production

COPY --from=api /compiler/application ./

EXPOSE 8080

CMD [ "./application" ]

FROM golang:1.23.0 AS build

WORKDIR /compiler

COPY go.mod go.sum /

RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o bin/order-service main.go

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /production

COPY --from=build /compiler/bin/order-service /app/

CMD [ "/app/order-service" ]

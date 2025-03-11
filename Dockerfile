FROM golang:1.24 AS build_stage

COPY . src/taskManager

WORKDIR /go/src/taskManager

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./app ./main.go

FROM scratch

WORKDIR /app

COPY --from=build_stage /go/src/taskManager/app .

EXPOSE 8080

ENTRYPOINT ["./app", "server"]
FROM golang:1.14.3-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /out/auth-service cmd/server.go

FROM scratch
COPY --from=build /out/auth-service /
COPY front /front
CMD ["/auth-service"]

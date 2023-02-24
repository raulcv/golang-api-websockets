ARG GO_VERSION=1.19.5

FROM golang:${GO_VERSION}-alpine as builder
RUN go env -w GOPROXY=direct
RUN apk add --no-cache git
RUN  apk --no-cache add ca-certificates && update-ca-certificates
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 go build \
  -installsuffix 'static' \
  -o /golang-api-ws
FROM scratch AS runner
COPY --from=builder/ect/ssl/cets/certificates.ca /ect/ssl/certs/
COPY .env ./
COPY --from=builder/golang-api-ws /golang-api-ws
EXPOSE 3020
ENTRYPOINT ["/golang-api-ws"]
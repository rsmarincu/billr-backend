FROM --platform=linux/amd64 golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /billr

## Deploy
FROM --platform=linux/amd64 gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /billr /billr
COPY internal/invoice_template.html internal/invoice_template.html

ARG DB_URL
ENV HTTP_PORT 7001
ENV ENVIRONMENT development
ARG GB_HOST

EXPOSE 7001


ENTRYPOINT ["/billr"]
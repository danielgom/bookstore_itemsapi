FROM golang:1.16.4-alpine as build

WORKDIR /application

COPY . .

RUN go build -o itemsapi

FROM alpine

# Envs to connect with elasticsearch host instance
ENV ELASTIC_SEARCH_HOST=host.docker.internal
ENV ELASTIC_SEARCH_PORTS=9200

WORKDIR /app

COPY --from=build /application/itemsapi /app/

EXPOSE 8080

CMD ["./itemsapi"]

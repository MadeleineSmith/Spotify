FROM golang:1.12.9-stretch

WORKDIR /
COPY ./ .

RUN go build -o spotify

EXPOSE 5000
CMD ["./spotify"]
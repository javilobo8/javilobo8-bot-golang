FROM golang:1.17.8-alpine

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . ./
RUN go build -o ./javilobo8-bot
CMD [ "./javilobo8-bot" ]
FROM golang:1.18-alpine

WORKDIR /src/backend

COPY . ./

RUN go mod download

RUN go build -o /ex_app

EXPOSE 3030

CMD [ "/ex_app", "--configPath", "empty.env"]
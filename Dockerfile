FROM golang:latest

WORKDIR /app

#COPY requirements.txt .
#RUN pip install -r requirements.txt

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8080

RUN go build -o main .

CMD ./main
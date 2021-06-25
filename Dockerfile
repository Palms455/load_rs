#FROM scratch
#ADD main /
#CMD ["/main"]


FROM golang:latest
RUN apt-get update && apt-get install -y build-essential libxml2-dev libxmlsec1-dev 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o main cmd/main.go 
CMD ["/app/main"]
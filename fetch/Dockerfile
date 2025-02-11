FROM golang:1.23.6
 
WORKDIR /fetch
COPY . .
RUN make build

CMD make run/live

FROM golang:latest as builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o monolith cmd/monolith/*.go


######## Start a new stage from scratch #######
FROM alpine:latest
RUN apk --no-cache add ca-certificates && update-ca-certificates
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/monolith .
COPY --from=builder /app/cmd/monolith/conf .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./monolith"]
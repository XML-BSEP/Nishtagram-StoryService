FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /src

# Copy go mod and sum files
COPY go.mod go.sum ./

COPY Nishtagram-Logger/ ../../Nishtagram-Logger/

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

RUN go get github.com/mattn/go-isatty@v0.0.12
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

######## Start a new stage from scratch #######
FROM alpine:latest  

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /src/config ./src/config
COPY --from=builder /src/certificate ./src/certificate
COPY --from=builder /src/main .
COPY --from=builder /src/assets ./src/assets
COPY --from=builder /src/domain ./src/domain
COPY --from=builder /src/dto ./src/dto
COPY --from=builder /src/gateway ./src/gateway
COPY --from=builder /src/http ./src/http
COPY --from=builder /src/infrastructure ./src/infrastructure
COPY --from=builder /src/interactor ./src/interactor
COPY --from=builder /src/repository ./src/repository
COPY --from=builder /src/usecase ./src/usecase


# Command to run the executable
CMD ["./main"]

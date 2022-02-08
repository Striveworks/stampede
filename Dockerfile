FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o stampede .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/stampede .

# Create nonroot user while we still have /bin/ash
RUN echo "nobody:x:65534:65534:nobody:/:/sbin/nologin" >> /tmp/nobody

# Build a small image
FROM scratch

COPY --from=builder /dist/stampede /

# Bring over the nobody user & su
COPY --from=builder /tmp/nobody /etc/passwd
USER nobody

# Command to run
ENTRYPOINT ["/main"]

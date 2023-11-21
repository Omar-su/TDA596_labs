# Use an official Go runtime as a parent image
FROM golang:1.16

# Set the working directory in the container
WORKDIR /

# Copy the entire project directory into the container at /app
COPY . .

RUN chmod +x server/server.go proxy/proxy.go

# Build the Go applications inside the container
RUN go build -o serverexe ./server/server.go
RUN go build -o proxyexe ./proxy/proxy.go

# Specify the commands to run your binaries
CMD ["./serverexe"]
# or CMD ["./proxy"] depending on which one you want to run by default
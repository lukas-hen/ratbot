FROM golang:1.18.0-stretch
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Build executable
RUN go build -o ./rattbot *.go

# Run the executable
CMD ["./entrypoint.sh"]
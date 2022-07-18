# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/mikejk8s/talesmud

WORKDIR /go/src/github.com/mikejk8s/talesmud

RUN go mod download
RUN go install github.com/mikejk8s/talesmud/cmd/tales

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/tales

# Document that the service listens on port 8080.
EXPOSE 8010
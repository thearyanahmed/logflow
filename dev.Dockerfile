FROM golang:1.16

RUN apt-get update

RUN apt-get install bash

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/thearyanahmed/logflow

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

COPY .env.example .env
# This container exposes port 5053 to the outside world
EXPOSE 5053

# Run the executable
#CMD ["logflow"]

ENTRYPOINT ["tail", "-f", "/dev/null"]

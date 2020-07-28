# Starting the MVP Studio base image, based on Debian distro.
FROM mvpstudio/base:v1 as go-build

# Download the Go distribution tar.gz.
RUN curl -o /usr/local/go1.14.4.linux-amd64.tar.gz https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz
# Decompress the tar and extract to the /usr/local directory.
RUN tar -C /usr/local -xzf /usr/local/go1.14.4.linux-amd64.tar.gz

# Set the Go path so that `go` program can be run from anywhere.
ENV PATH="${PATH}:/usr/local/go/bin"

# Change the working directory inside the docker image.
WORKDIR /home/go/src/
# Copy the go files from local machine to the docker image.
COPY *.go ./
COPY go.* ./
# Build the Go binary.
RUN go build -o build/goserver

# This is a multi-stage build. So start fresh with another base image.
FROM mvpstudio/base:v1
WORKDIR /home/mvp/app
# This time, copy over the binary that was built on the first base image.
COPY --from=go-build /home/go/src/build .
COPY frontend-template.html .

# It's best practice to run as a user, not as root.
USER mvp
# Run the go app.
ENTRYPOINT ["./goserver"]

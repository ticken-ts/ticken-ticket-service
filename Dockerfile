# syntax=docker/dockerfile:
# specify the base image to be used for the application, alpine or ubuntu
FROM golang:1.18-alpine as build

# git is required to fetch go dependencies
RUN apk add --no-cache ca-certificates git

# add file with credetentials to
# download private go modules
COPY .netrc /root/.netrc
RUN chmod 600 /root/.netrc

# create a working directory inside the image
WORKDIR /src

# copy Go modules and dependencies to image
COPY go.mod ./

ENV GOPRIVATE="github.com/ticken-ts/ticken-pvtbc-connector"
RUN go mod download

# copy directory files i.e all files ending with .go
COPY . .

# compile application
RUN CGO_ENABLED=0 go build -o /service .

FROM scratch AS final

COPY --from=build /service /service

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 9000

# command to be used to execute when the image is used to start a container
ENTRYPOINT [ "/service" ]


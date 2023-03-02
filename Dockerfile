ARG GO_V=1.19
ARG ALPINE_V=3.16

FROM golang:${GO_V}-alpine${ALPINE_V} AS builder

# Sets PWD
WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the codebase into the layer's pwd. Check .dockerignore if you're missing something.
COPY . .

# Builds an executable named main.
RUN go build ./cmd/app/main.go

# Use a thin layer to run the software.
FROM alpine:${ALPINE_V}

# Place the executable named main where any user may execute it: /usr/local/bin/
WORKDIR /usr/local/bin
COPY --from=builder /app/main .

# The ENTRYPOINT specifies a command that will always be executed when the container starts.
# Docker has a default entrypoint which is /bin/sh -c but does not have a default command.
# The CMD specifies arguments that will be fed to the ENTRYPOINT.
# TL;DR: We execute ./main when starting the container.
CMD ["./main"]
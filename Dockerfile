ARG GOLANG_VER=latest
ARG ALPINE_VER=latest

FROM golang:${GOLANG_VER} AS builder
WORKDIR /go/src/app
COPY go.* *.go ./
COPY cmd cmd/
ENV CGO_ENABLED=0
ARG ACTIONLINT_VER=
RUN go build -v -ldflags "-s -w -X github.com/zew1me/actionlint.version=${ACTIONLINT_VER}" ./cmd/actionlint

FROM koalaman/shellcheck-alpine:stable AS shellcheck

FROM alpine:${ALPINE_VER}
LABEL org.opencontainers.image.source="https://github.com/zew1me/actionlint" \
      org.opencontainers.image.title="actionlint-zew1me" \
      org.opencontainers.image.description="Fork of rhysd/actionlint: static checker for GitHub Actions workflow files"
COPY --from=builder /go/src/app/actionlint /usr/local/bin/
COPY --from=shellcheck /bin/shellcheck /usr/local/bin/shellcheck
RUN apk add --no-cache py3-pyflakes
# Numeric UID:GID of Alpine's built-in `guest` user (405:100) so it resolves on any host (hadolint DL3066)
USER 405:100
ENTRYPOINT ["/usr/local/bin/actionlint"]

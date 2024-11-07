FROM docker.io/library/golang:1.23.3-alpine3.19 AS build

RUN apk add git make

WORKDIR /workspace
ADD ./ /workspace

RUN make -C cmd/autharr install \
      DESTDIR=/cache \
      PREFIX=/usr \
      VERSION=${VERSION}

RUN make -C cmd/healarr install \
      DESTDIR=/cache \
      PREFIX=/usr \
      VERSION=${VERSION}

FROM docker.io/library/alpine:3.20

COPY --from=build /cache /
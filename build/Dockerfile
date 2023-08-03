FROM golang:latest as builder

WORKDIR /build
COPY . .

RUN make all


FROM ubuntu:latest

WORKDIR /

COPY --from=builder /build/bin/secretary         bin/secretary
COPY --from=builder /build/configs               configs
COPY --from=builder /build/secrets               secrets
COPY --from=builder /build/Makefile              .

RUN chmod +x /bin/secretary
RUN apt-get update
RUN apt-get install -y make ca-certificates openssl

ENTRYPOINT ["make"]
CMD ["run"]

FROM alpine:latest as builder
ARG TARGETOS
ARG TARGETARCH
COPY ./build/fiber-gpt-${TARGETOS}-${TARGETARCH} .
RUN mv fiber-gpt-${TARGETOS}-${TARGETARCH} server
RUN chmod +x server

FROM alpine:latest as runtime
WORKDIR /app
COPY --from=builder /server .

ENV FIBER_GPT_SERVER_PORT 3000
ENV FIBER_GPT_CONFIG_PATH /var/fiber-gpt/config.yaml
ENV FIBER_GPT_DB_URL /var/fiber-gpt/db.sqlite

EXPOSE 3000
ENTRYPOINT ["/app/server"]
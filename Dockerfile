FROM ubuntu:latest

EXPOSE 8080

WORKDIR /app

RUN useradd senketsu
USER senketsu

COPY senketsu .
COPY templates ./templates

CMD ["./senketsu"]

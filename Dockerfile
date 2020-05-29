FROM hayd/deno:latest

EXPOSE 8080

WORKDIR /app
USER deno

COPY . /app
RUN deno cache index.ts

CMD ["run", "--allow-net", "index.ts"]

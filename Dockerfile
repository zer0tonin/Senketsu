FROM hayd/deno:latest

EXPOSE 8080

WORKDIR /app
USER deno

COPY . /app
RUN deno cache src/app.ts

CMD ["run", "--allow-net", "--allow-read", "src/app.ts"]

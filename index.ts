import { serve } from 'https://deno.land/std/http/server.ts';

const server = serve({ port: 8080 });
console.log("starting server on localhost:8080")
for await (const req of server) {
  req.respond({ body: "Hello World\n" });
}

import { Application, Router, Context } from "https://deno.land/x/oak/mod.ts";

const router = new Router();
router
  .get("/", (context) => {
    context.response.body = "Hello World!";
  })
  .get("/:string", (context) => {
    context.response.body = "Hello " + context.params.string;
  });

const app = new Application();
app.use(router.routes());
app.use(router.allowedMethods());

console.log("Starting server on port 8080");
await app.listen({ port: 8080 });

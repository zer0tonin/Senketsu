import { Router } from "https://deno.land/x/oak/mod.ts";

const router = new Router();
router
  .get("/", (context) => {
    context.response.body = "Hello World!";
  })
  .get("/:string", async (context) => {
    context.render('templates/index.ejs', { name: context.params.string });
  });

export default router;

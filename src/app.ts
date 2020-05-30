import * as log from "https://deno.land/std/log/mod.ts";
import { Application } from "https://deno.land/x/oak/mod.ts";
import { organ } from "https://raw.githubusercontent.com/denjucks/organ/master/mod.ts";
import { Snelm } from "https://deno.land/x/snelm/mod.ts";
import {
  viewEngine,
  engineFactory,
  adapterFactory,
} from "https://deno.land/x/view_engine/mod.ts";

import router from "./router.ts";

const app = new Application();

// request logging
app.use(organ());

// error logging
app.use(async (ctx, next) => {
  try {
    await next();
  } catch(err) {
    log.error(err)
    throw err;
  }
});

// templating
const ejsEngine = engineFactory.getEjsEngine();
const oakAdapter = adapterFactory.getOakAdapter();
app.use(viewEngine(oakAdapter, ejsEngine));

// security headers
const snelm = new Snelm("oak");
await snelm.init();
app.use((ctx, next) => {
  ctx.response = snelm.snelm(ctx.request, ctx.response);
  
  next();
});

// routing
app.use(router.routes());
app.use(router.allowedMethods());

log.info("Starting server on port 8080");
await app.listen({ port: 8080 });

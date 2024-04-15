import Elysia from "elysia";
import corsConfig from "./config/cors";
import usersController from "./controllers/users-controller";

const port = process.env.PORT || 3000;
new Elysia({
  cookie: {
    secrets: process.env.COOKIE_AUTH_SECRET,
    sign: ["auth"],
  },
})
  .use(corsConfig())
  .use(usersController)
  .onError(() => {
    return new Response("INTERNAL SERVER ERROR", { status: 500 });
  })
  .listen(port, () => console.log("Server is listening at port", port));

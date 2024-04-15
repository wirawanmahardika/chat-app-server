import Elysia from "elysia";

const port = process.env.PORT || 3000;
new Elysia()
  .post("/", () => {})
  .listen(port, () => console.log("Server is listening at port", port));

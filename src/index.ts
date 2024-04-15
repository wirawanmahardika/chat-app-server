import Elysia, { t } from "elysia";
import prisma from "./app/prisma";
import { toBuffer } from "bun:ffi";
import corsConfig from "./config/cors";

const port = process.env.PORT || 3000;
new Elysia()
  .use(corsConfig())
  .post(
    "/signup",
    async ({ body }) => {
      const countUser = await prisma.users.count({
        where: { username: body.username },
      });

      if (countUser > 0)
        return new Response("username is not available", { status: 400 });

      await prisma.users.create({
        data: {
          username: body.username,
          fullname: body.fullname,
          email: body.email,
          password: body.password,
          photo_profile: Buffer.from(await body.photo_profile.arrayBuffer()),
        },
      });

      return new Response("Berhasil signup", { status: 201 });
    },
    {
      type: "multipart/form-data",
      body: t.Object({
        fullname: t.String({ error: "fullname should not be empty" }),
        email: t.String({ format: "email", error: "email format is invalid" }),
        username: t.String({
          minLength: 6,
          error: "username should have at least 6 characters",
        }),
        password: t.String({
          minLength: 6,
          error: "password should have at least 6 characters",
        }),
        photo_profile: t.File({
          type: ["image/jpeg", "image/png"],
          maxSize: 5_242_880,
          error: "invalid file",
        }),
      }),
    }
  )
  .post("/", () => {}, {
    body: t.Object({
      username: t.String({ minLength: 6 }),
      password: t.String({ minLength: 6 }),
    }),
  })
  .listen(port, () => console.log("Server is listening at port", port));

import Elysia, { Cookie } from "elysia";
import usersServices from "../services/users-services";
import jwtConf from "../config/jwt";

const usersRoute = new Elysia({ prefix: "/users" })
  .use(jwtConf())
  .post(
    "/signup",
    async ({ body }) => {
      const countUser = await usersServices.signup.repository.countUser(
        body.username
      );

      if (countUser > 0)
        return new Response("username is not available", { status: 400 });

      const data = {
        username: body.username,
        fullname: body.fullname,
        email: body.email,
        password: body.password,
        photo_profile: Buffer.from(await body.photo_profile.arrayBuffer()),
      };

      await usersServices.signup.repository.createUser(data);
      return new Response("Berhasil signup", { status: 201 });
    },
    usersServices.signup.schema
  )
  .post(
    "/login",
    async ({ body, cookie, jwt }) => {
      const user = await usersServices.login.repository.getUser(body.username);

      if (!user)
        return new Response("username tidak pernah terdaftar", { status: 401 });

      if (user.password !== body.password) {
        return new Response("password salah", { status: 401 });
      }

      cookie.auth.set({
        httpOnly: true,
        maxAge: 3600 * 24,
        path: "/",
        priority: "high",
        secrets: process.env.COOKIE_AUTH_SECRET || "asdfoaur8eqw795873948",
        value: await jwt.sign({ username: body.username }),
        sameSite: "strict",
        secure: process.env.IS_HTTPS ? true : false,
      });

      return "Berhasil login";
    },
    usersServices.login.schema
  );

const userRoute = new Elysia({ prefix: "/user" })
  .use(jwtConf())
  .derive(async ({ jwt, cookie }) => {
    const jwtData = await jwt.verify(cookie.auth.value);
    return {
      authenticated: jwtData ? true : false,
      user: jwtData ? jwtData : { username: "" },
    };
  })
  .onBeforeHandle(async ({ authenticated }) => {
    if (!authenticated) {
      return new Response("Membutuhkan login terlebih dahulu", { status: 401 });
    }
  })
  .get("/info", async ({ user: { username } }) => {
    const userData = await usersServices.info.repository.getUserData(username);
    return { ...userData, photo_profile: "" };
  });

export default new Elysia({ prefix: "/api/v1" }).use(usersRoute).use(userRoute);

import Elysia from "elysia";
import usersServices from "../services/users-services";

export default new Elysia({ prefix: "/api/v1/users" })
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
    async ({ body }) => {
      const user = await usersServices.login.repository.getUser(body.username);

      if (!user)
        return new Response("username tidak pernah terdaftar", { status: 401 });

      if (user.password !== body.password) {
        return new Response("password salah", { status: 401 });
      }

      return "Berhasil login";
    },
    usersServices.login.schema
  );

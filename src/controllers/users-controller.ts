import Elysia from "elysia";
import usersServices from "../services/users-services";

export default new Elysia({ prefix: "/api/v1/users" }).post(
  "/signup",
  async ({ body }) => {
    const countUser = await usersServices.signup.signupRepositories.countUser(
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

    await usersServices.signup.signupRepositories.createUser(data);
    return new Response("Berhasil signup", { status: 201 });
  },
  usersServices.signup.signupBodyValidation
);

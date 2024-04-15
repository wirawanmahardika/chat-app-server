import Elysia from "elysia";
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
        id: undefined,
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
        value: await jwt.sign({ username: body.username, id: user.id }),
        sameSite: "strict",
        expires: new Date(Date.now() + 1000 * 3600 * 24),
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
      user: jwtData ? jwtData : { username: "", id: "" },
    };
  })
  .onBeforeHandle(async ({ authenticated }) => {
    if (!authenticated) {
      return new Response("Membutuhkan login terlebih dahulu", { status: 401 });
    }
  })
  .get("/info", async ({ user: { username, id } }) => {
    const userData = await usersServices.info.repository.getUserData(username);
    return {
      ...userData,
      photo_profile: process.env.SERVER_URL + "/api/v1/user/photo/" + id,
    };
  })
  .get(
    "/photo/:id_user",
    async ({ params }) => {
      const photo = await usersServices.photo.repository.getPhoto(
        params.id_user
      );
      if (photo) {
        return new Response(photo);
      } else {
        return new Response("photo doesnt exist", { status: 404 });
      }
    },
    usersServices.photo.schema
  )
  .delete("/", ({ cookie }) => {
    cookie.auth.remove();
    return `Berhasil logout`;
  })
  .post(
    "/add-friend",
    async ({ body, user }) => {
      if (user.id === body.id_friend) {
        return new Response(
          "Tidak bisa menjalin pertemanan dengan diri sendiri",
          { status: 500 }
        );
      }

      const response = await usersServices.addFriend.repository.addFriend(
        user.id,
        body.id_friend
      );

      return new Response(response);
    },
    usersServices.addFriend.schema
  );
// .get("/friend-requests", async ({ user }) => {
//   const requests =
//     await usersServices.friendRequests.repository.getAllRequests(user.id);
// });

export default new Elysia({ prefix: "/api/v1" }).use(usersRoute).use(userRoute);

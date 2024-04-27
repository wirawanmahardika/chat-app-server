import Elysia from "elysia";
import jwtConf from "../config/jwt";
import friendServices from "../services/friend-services";

export default new Elysia({ prefix: "/api/v1/friend" })
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
  .post(
    "/",
    async ({ body, user }) => {
      if (user.id === body.id_friend) {
        return new Response(
          "Tidak bisa menjalin pertemanan dengan diri sendiri",
          { status: 500 }
        );
      }

      const response = await friendServices.addFriend.repository.addFriend(
        user.id,
        body.id_friend
      );

      return new Response(response);
    },
    friendServices.addFriend.schema
  )
  .get("/relationship-requests", async ({ user }) => {
    const requests =
      await friendServices.friendRequests.repository.getAllRequests(user.id);

    return requests.map((r) => {
      return {
        id_friendship: r.id_friendship,
        photo_profile:
          process.env.SERVER_URL + "/api/v1/user/photo/" + r.user_1.id,
        fullname: r.user_1.fullname,
        created_at: r.created_at,
        status: r.status,
      };
    });
  })
  .patch(
    "/response-to-request",
    async ({ body }) => {
      const status =
        await friendServices.requestResponse.repository.updateFriendshipStatus(
          body.id_friendship,
          body.status,
          body.rejection
        );

      if (status === "friend") {
        return "Berhasil menjalin pertemanan";
      } else {
        return "Berhasil menghapus permintaan pertemenan";
      }
    },
    friendServices.requestResponse.schema
  )
  .get("/get-all", async ({ user }) => {
    const friends = await friendServices.friends.repository.getFriends(user.id);

    return friends.map((f) => {
      return {
        ...f,
        photo_profile: process.env.SERVER_URL + "/api/v1/user/photo/" + f.id,
      };
    });
  })
  .get(
    "/status/:id_friend",
    async ({ params }) => {
      const friendStatus =
        await friendServices.friendStatus.repository.getFriendStatus(
          params.id_friend
        );
      return friendStatus;
    },
    friendServices.friendStatus.schema
  );

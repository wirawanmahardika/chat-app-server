import Elysia, { t } from "elysia";
import jwtConf from "../config/jwt";
import { JWTPayloadSpec } from "@elysiajs/jwt";
import {
  createMessage,
  getChatMessages,
  getFriendStatus,
  getLastMessageOfEachFriend,
  updateUserStatus,
} from "../services/chat-services";
import Stream from "@elysiajs/stream";
import prisma from "../app/prisma";

type jwtPayloadSchema = {
  username: string;
  id: string;
} & JWTPayloadSpec;

export default new Elysia()
  .use(jwtConf())
  .derive(async ({ jwt, cookie }) => {
    const jwtData = await jwt.verify(cookie.auth.value);
    return {
      authenticated: jwtData ? true : false,
      user: jwtData ? jwtData : { username: "", id: "" },
    };
  })
  .ws("/ws", {
    body: t.Object({
      type: t.Union([
        t.Literal("subscribe"),
        t.Literal("chat"),
        t.Literal("leave"),
      ]),
      data: t.Any(),
    }),
    async beforeHandle({ cookie, jwt }) {
      const result = await jwt.verify(cookie.auth.value);
      if (!result) {
        return "Failed to connect to websocket";
      }
    },

    async message(ws, { type, data }) {
      const dataFromToken = (await ws.data.jwt.verify(
        ws.data.cookie.auth.value
      )) as jwtPayloadSchema;
      ws.id = dataFromToken.id;

      if (type === "subscribe") {
        ws.subscribe(data.channel);
        await updateUserStatus(dataFromToken.id, true);
        ws.publish(
          data.channel,
          JSON.stringify({ type: "join", status: true })
        );
      }

      if (type === "chat") {
        await createMessage({
          from: dataFromToken.id,
          to: data.to,
          message: data.message,
          id_friendship: data.channel,
        });

        ws.publish(
          data.channel,
          JSON.stringify({
            type: "chat",
            data: { message: data.message, from: dataFromToken.username },
          })
        );
      }

      if (type === "leave") {
        ws.publish(
          data.channel,
          JSON.stringify({
            type: "leave",
            data: { status_friend: false },
          })
        );
        ws.unsubscribe(data.channel);
        ws.terminate();
      }
    },

    async close(ws) {
      await updateUserStatus(ws.data.user.id, false);
    },
  })
  .onBeforeHandle(async ({ jwt, cookie }) => {
    const result = await jwt.verify(cookie.auth.value);

    if (!result) {
      return new Response("Anda perlu login terlebih dahulu", { status: 401 });
    }
  })
  .get(
    "/api/v1/chats/:id_friendship",
    async ({ params }) => {
      const messages = await getChatMessages(params.id_friendship);
      return messages;
    },
    { params: t.Object({ id_friendship: t.String() }) }
  )
  .get(
    "/sse",
    ({ user }) =>
      new Stream((stream) => {
        setInterval(async () => {
          const friendsMessages = await getLastMessageOfEachFriend(user.id);
          stream.send(friendsMessages);
        }, 1200);
      })
  )
  .get(
    "/api/v1/chats/friend_status/:id_friend",
    async ({ params }) => {
      const friendStatus = await getFriendStatus(params.id_friend);
      return friendStatus;
    },
    { params: t.Object({ id_friend: t.String({ format: "uuid" }) }) }
  );

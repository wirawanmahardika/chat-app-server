import { t } from "elysia";
import prisma from "../app/prisma";

const addFriendSchema = {
  body: t.Object({
    id_friend: t.String({
      format: "uuid",
      error: "invalid ID, tolong cek kembali",
    }),
  }),
};

class addFriendRepository {
  static async addFriend(id_user_1: string, id_user_2: string) {
    const friendship = await prisma.friendships.findFirst({
      where: {
        OR: [
          { id_user_1: id_user_1, id_user_2: id_user_2 },
          { id_user_1: id_user_2, id_user_2: id_user_1 },
        ],
      },
      select: { id_friendship: true },
    });

    if (friendship) {
      await prisma.friendships.update({
        where: { id_friendship: friendship.id_friendship },
        data: {
          status: "friend",
        },
      });

      return "Berhasil menjalin pertemanan";
    }

    await prisma.friendships.create({
      data: {
        id_user_1: id_user_1,
        id_user_2: id_user_2,
        status: "pending",
      },
    });

    return "Berhasil mengirim permintaan pertemanan";
  }
}

class friendRequestsRepository {
  static async getAllRequests(id_user: string) {
    return prisma.friendships.findMany({
      where: {
        id_user_2: id_user,
        status: { notIn: ["friend", "blocked"] },
      },
      include: {
        user_1: {
          select: {
            email: true,
            id: true,
            fullname: true,
            username: true,
          },
        },
      },
    });
  }
}

const requestResponseSchema = {
  body: t.Object({
    id_friendship: t.String({ format: "uuid", error: "invalid id" }),
    rejection: t.Optional(t.Boolean({ default: false })),
    status: t.Optional(
      t.Union(
        [t.Literal("pending"), t.Literal("friend"), t.Literal("blocked")],
        {
          error: "allowed status are pending, friend, or block",
          default: "pending",
        }
      )
    ),
  }),
};

class requestResponseRepository {
  static async updateFriendshipStatus(
    id_friendship: string,
    status: "pending" | "friend" | "blocked" | undefined,
    rejection: boolean | undefined
  ) {
    if (rejection) {
      await prisma.friendships.delete({ where: { id_friendship } });
      return null;
    } else {
      const res = await prisma.friendships.update({
        where: { id_friendship },
        data: {
          status: status,
        },
        select: {
          status: true,
        },
      });
      return res.status;
    }
  }
}

class friendsRepository {
  static async getFriends(id_user: string) {
    const usersByUser1 = await prisma.friendships.findMany({
      where: {
        id_user_1: id_user,
        status: "friend",
      },
      include: {
        user_2: {
          select: {
            username: true,
            fullname: true,
            id: true,
            email: true,
          },
        },
      },
    });
    const usersByUser2 = await prisma.friendships.findMany({
      where: {
        id_user_2: id_user,
        status: "friend",
      },
      include: {
        user_1: {
          select: {
            username: true,
            fullname: true,
            id: true,
            email: true,
          },
        },
      },
    });

    const friends: {
      id: string;
      fullname: string;
      username: string;
      email: string;
      id_friendship: string;
      message?: string;
      from?: string;
      created_at?: Date;
    }[] = [
      ...usersByUser1.map((u) => {
        return {
          id_friendship: u.id_friendship,
          ...u.user_2,
        };
      }),
      ...usersByUser2.map((u) => {
        return {
          id_friendship: u.id_friendship,
          ...u.user_1,
        };
      }),
    ];

    for (const f of friends) {
      const chat = await prisma.chat.findFirst({
        where: { id_friendship: f.id_friendship },
        orderBy: { created_at: "desc" },
        select: { message: true, from: true, created_at: true },
      });

      f.message = chat?.message;
      f.from = chat?.from;
      f.created_at = chat?.created_at;
    }

    return friends;
  }
}

const friendStatusSchema = {
  params: t.Object({ id_friend: t.String({ format: "uuid" }) }),
};

class friendStatusRepository {
  static async getFriendStatus(id_friend: string) {
    const result = await prisma.users.findUnique({
      where: { id: id_friend },
      select: {
        status: true,
      },
    });

    return result?.status;
  }
}

export default {
  addFriend: {
    repository: addFriendRepository,
    schema: addFriendSchema,
  },
  friendRequests: {
    repository: friendRequestsRepository,
  },
  requestResponse: {
    repository: requestResponseRepository,
    schema: requestResponseSchema,
  },
  friends: {
    repository: friendsRepository,
  },
  friendStatus: {
    repository: friendStatusRepository,
    schema: friendStatusSchema,
  },
};

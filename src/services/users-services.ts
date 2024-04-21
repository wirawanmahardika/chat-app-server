import { t } from "elysia";
import prisma from "../app/prisma";
import { usersType } from "../types/users-type";

const signupSchema = {
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
};

class signupRepositories {
  static async countUser(username: string) {
    return prisma.users.count({ where: { username: username } });
  }

  static async createUser(user: usersType) {
    return prisma.users.create({ data: user });
  }
}

const loginSchema = {
  body: t.Object({
    username: t.String({ error: "username is required" }),
    password: t.String({ error: "password is required" }),
  }),
};

class loginRepositories {
  static async getUser(username: string) {
    return prisma.users.findUnique({ where: { username } });
  }
}

class infoRepositories {
  static async getUserData(username: string) {
    return prisma.users.findUnique({
      where: { username },
      select: {
        id: true,
        email: true,
        username: true,
        fullname: true,
      },
    });
  }
}

const photoSchema = {
  params: t.Object({
    id_user: t.String(),
  }),
};

class photoRepositories {
  static async getPhoto(id_user: string) {
    const data = await prisma.users.findUnique({
      where: { id: id_user },
      select: { photo_profile: true },
    });
    return data?.photo_profile;
  }
}

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

class friendRequests {
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

export default {
  signup: {
    schema: signupSchema,
    repository: signupRepositories,
  },
  login: {
    schema: loginSchema,
    repository: loginRepositories,
  },
  info: {
    repository: infoRepositories,
  },
  photo: {
    schema: photoSchema,
    repository: photoRepositories,
  },
  addFriend: {
    schema: addFriendSchema,
    repository: addFriendRepository,
  },
  friendRequests: {
    repository: friendRequests,
  },
  requestResponse: {
    schema: requestResponseSchema,
    repository: requestResponseRepository,
  },
  friends: {
    repository: friendsRepository,
  },
};

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
};

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

export default {
  signup: {
    schema: signupSchema,
    repository: signupRepositories,
  },
  login: {
    schema: loginSchema,
    repository: loginRepositories,
  },
};

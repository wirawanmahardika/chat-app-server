import { t } from "elysia";
import prisma from "../app/prisma";
import { usersType } from "../types/users-type";

export const signupBodyValidation = {
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

export class signupRepositories {
  static async countUser(username: string) {
    return prisma.users.count({ where: { username: username } });
  }

  static async createUser(user: usersType) {
    return prisma.users.create({ data: user });
  }
}

export default {
  signup: {
    signupBodyValidation,
    signupRepositories,
  },
};

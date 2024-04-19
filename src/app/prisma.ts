import { PrismaClient } from "@prisma/client";

const prisma = new PrismaClient({
  // log: [
  //   {
  //     emit: "stdout",
  //     level: "query",
  //   },
  // ],
});

export default prisma;

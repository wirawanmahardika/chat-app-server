import prisma from "../app/prisma";

export async function getChatMessages(id_friendship: string) {
  return prisma.chat.findMany({
    where: { id_friendship },
    orderBy: { created_at: "asc" },
    take: 25,
  });
}

export async function createMessage(data: {
  id_friendship: string;
  from: string;
  to: string;
  message: string;
}) {
  return prisma.chat.create({ data });
}

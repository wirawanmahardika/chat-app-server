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

export async function updateUserStatus(id: string, status: boolean) {
  await prisma.users.update({ where: { id }, data: { status: status } });
}

export async function getLastMessageOfEachFriend(id: string) {
  const usersByUser1 = await prisma.friendships.findMany({
    where: {
      id_user_1: id,
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
      id_user_2: id,
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
    photo_profile?: string;
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
    f.photo_profile = process.env.SERVER_URL + "/api/v1/user/photo/" + f.id;
  }

  return friends;
}

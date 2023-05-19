import { Prisma, PrismaClient } from '@prisma/client';

const prisma = new PrismaClient();

export async function getTasksWithTagsOfUser(user_id: number) {
    const tasks = await prisma.tasks.findMany({
        include: { Tags: true },
        where: { user_id: 1 },
    });

    return tasks;
}

export async function getTagsOfUser(user_id: number) {
    const tags = await prisma.tags.findMany({
        where: { user_id: 1 },
    });

    return tags;
}

export type TaskWithTags = Prisma.PromiseReturnType<
    typeof getTasksWithTagsOfUser
>[number];

export type Tag = Prisma.PromiseReturnType<typeof getTagsOfUser>[number];

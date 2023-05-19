import { Prisma, PrismaClient } from '@prisma/client';

declare global {
    var prisma: PrismaClient | undefined;
}

export const prisma = global.prisma || new PrismaClient();

if (process.env.NODE_ENV === 'development') global.prisma = prisma;

export async function getTasksWithTagsOfUser(user_id: number) {
    const tasks = await prisma.tasks.findMany({
        include: { Tags: true },
        where: { user_id },
    });

    return tasks;
}

export async function getTagsOfUser(user_id: number) {
    const tags = await prisma.tags.findMany({
        where: { user_id },
    });

    return tags;
}

export type TaskWithTags = Prisma.PromiseReturnType<
    typeof getTasksWithTagsOfUser
>[number];

export type Tag = Prisma.PromiseReturnType<typeof getTagsOfUser>[number];

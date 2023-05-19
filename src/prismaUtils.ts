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

export async function getTagsWithTaskCountOfUser(user_id: number) {
    const tags = await prisma.tags.findMany({
        where: { user_id },
        include: { _count: { select: { tasks: true } } },
    });

    // No comfortable solution to just add the field and have typescript resolve the type
    return tags.map(({ _count, ...tag }) => ({
        task_count: _count.tasks,
        ...tag,
    }));
}

export type TaskWithTags = Prisma.PromiseReturnType<
    typeof getTasksWithTagsOfUser
>[number];

export type Tag = Prisma.PromiseReturnType<typeof getTagsOfUser>[number];
export type TagWithTaskCount = Prisma.PromiseReturnType<
    typeof getTagsWithTaskCountOfUser
>[number];

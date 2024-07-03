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

export async function getTodoLists(user_id: number) {
    const todoLists = await prisma.todoLists.findMany({
        where: { user_id },
        include: { _count: { select: { Tasks: true } } },
    });

    // No comfortable solution to just add the field and have typescript resolve the type
    return todoLists.map(({ _count, ...todoList }) => ({
        task_count: _count.Tasks,
        ...todoList,
    }));
}

export async function getTodoListsWithTasks(user_id: number) {
    const todoLists = await prisma.todoLists.findMany({
        where: { user_id },
        include: { Tasks: true },
    });

    return todoLists;
}

export async function getTodoListWithTasksById(todolist_id: number) {
    const todoList = await prisma.todoLists.findFirst({
        where: { id: todolist_id },
        include: { Tasks: { include: { Task: true } } },
    });

    return {
        ...todoList,
        Tasks: todoList?.Tasks.map(({ Task, ...t }) => ({ ...t, ...Task })),
    };
}

export type TaskWithTags = Prisma.PromiseReturnType<
    typeof getTasksWithTagsOfUser
>[number];

export type Tag = Prisma.PromiseReturnType<typeof getTagsOfUser>[number];
export type TagWithTaskCount = Prisma.PromiseReturnType<
    typeof getTagsWithTaskCountOfUser
>[number];

export type TodoList = Prisma.PromiseReturnType<typeof getTodoLists>[number];
export type TodoListWithTasks = Prisma.PromiseReturnType<
    typeof getTodoListsWithTasks
>[number];

export type TodoListWithTasksData = Prisma.PromiseReturnType<
    typeof getTodoListWithTasksById
>;

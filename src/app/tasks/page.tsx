import TaskForm from '@/components/TaskForm';

import {
    TaskWithTags,
    getTagsOfUser,
    getTasksWithTagsOfUser,
} from '@/prismaUtils';
import { Tags } from '@prisma/client';
import { TaskList } from './TasksList';

const user_id = 1;

async function getTasks() {
    try {
        const tasks = await getTasksWithTagsOfUser(user_id);

        return { tasks };
    } catch (error) {
        return { error, tasks: [] as TaskWithTags[] };
    }
}

async function getTags() {
    try {
        const tags = await getTagsOfUser(user_id);

        return { tags };
    } catch (error) {
        return { error, tags: [] as Tags[] };
    }
}

export default async function Page() {
    const { tasks, error: taskError } = await getTasks();
    const { tags, error: tagError } = await getTags();

    if (taskError) {
        console.error(taskError);
    }

    if (tagError) {
        console.error(tagError);
    }

    return (
        <div>
            <h1 className="text-2xl">Tasks</h1>
            <TaskForm tags={tags} />

            <br />

            <TaskList tasks={tasks} tags={tags} />
        </div>
    );
}

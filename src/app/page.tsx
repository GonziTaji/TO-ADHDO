import TaskForm from '@/components/TaskForm';
import { TaskList } from '@/components/TasksList';
import { Tags } from '@prisma/client';
import {
    TaskWithTags,
    getTagsOfUser,
    getTasksWithTagsOfUser,
} from '@/prismaUtils';

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

export default async function Home() {
    const { tasks, error: taskError } = await getTasks();
    const { tags, error: tagError } = await getTags();

    if (taskError) {
        console.error(taskError);
    }

    if (tagError) {
        console.error(tagError);
    }

    return (
        <main>
            <TaskForm tags={tags} />

            <br />

            <TaskList tasks={tasks} />
        </main>
    );
}

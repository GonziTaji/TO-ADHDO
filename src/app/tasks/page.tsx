import {
    TaskWithTags,
    getTagsOfUser,
    getTasksWithTagsOfUser,
} from '@/prismaUtils';
import { Tags } from '@prisma/client';
import NewTaskForm from './NewTaskForm';
import TaskItemList from './TaskItemList';

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
            <h1 className="text-2xl">Task List</h1>

            <section className="flex flex-col">
                <div className="border border-slate-400 p-4 my-2">
                    <NewTaskForm tags={tags} />
                </div>

                {tasks.map((task, i) => (
                    <div key={i} className="border-b border-slate-500 py-2">
                        <TaskItemList task={task} tags={tags} />
                    </div>
                ))}
            </section>
        </div>
    );
}

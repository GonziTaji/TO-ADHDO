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
            <h1 className="text-2xl my-2">Task List</h1>

            <section className="flex flex-col gap-6">
                <div className="border border-slate-400 p-4">
                    <NewTaskForm tags={tags} />
                </div>

                <div className="border border-gray-400">
                    {tasks.map((task, i) => (
                        <div
                            key={i}
                            className="border-b last:border-0 border-slate-500 p-4"
                        >
                            <TaskItemList task={task} tags={tags} />
                        </div>
                    ))}
                </div>
            </section>
        </div>
    );
}

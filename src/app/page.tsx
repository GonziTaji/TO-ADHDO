import TaskForm from '@/components/TaskForm';
import { TaskList } from '@/components/TasksList';
import { Task } from '@/types';
import { PrismaClient } from '@prisma/client';
console.log('in file');
async function getTasks() {
    console.log('in get tasks');
    const prisma = new PrismaClient();

    try {
        const tasks = await prisma.tasks.findMany({
            include: { Tags: true },
            where: { user_id: 1 },
        });

        console.log(tasks);

        return { tasks };
    } catch (error) {
        return { error, tasks: [] as Task[] };
    }
}

async function getTags() {
    const prisma = new PrismaClient();

    try {
        const tags: any[] = await prisma.tags.findMany({
            where: { user_id: 1 },
        });

        console.log(tags);

        return { tags };
    } catch (error) {
        return { error, tags: [] as any[] };
    }
}

export default async function Home() {
    return <p>i work</p>;
    const { tasks, error: taskError } = await getTasks();
    // const { tags, error: tagError } = await getTags();

    if (taskError) {
        console.error(taskError);
    }

    // if (tagError) {
    //     console.error(tagError);
    // }

    function addTask(newTask: Task) {
        // setTasks([...tasks, newTask]);
        console.log(newTask);
    }

    return (
        <main className="max-w-md mx-auto">
            <TaskForm tags={[]} />

            <br />

            <TaskList tasks={tasks as any} />
        </main>
    );
}

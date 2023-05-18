import TaskForm from '@/components/TaskForm';
import { TaskList } from '@/components/TasksList';
import { Task } from '@/types';
import { useState } from 'react';

export default async function Home() {
    const [tasks, setTasks] = useState<Task[]>([]);

    function addTask(newTask: Task) {
        setTasks([...tasks, newTask]);
    }

    return <p>hehe</p>;

    return (
        <main className="max-w-md mx-auto">
            <TaskForm addTask={addTask} />

            <br />

            <TaskList tasks={tasks} />
        </main>
    );
}

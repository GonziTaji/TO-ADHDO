import TaskForm from '@/components/TaskForm';
import { TaskList } from '@/components/TasksList';
import { Task } from '@/types';
import { useState } from 'react';

export default function Home() {
    const [tasks, setTasks] = useState<Task[]>([]);

    function addTask(newTask: Task) {
        setTasks([...tasks, newTask]);
    }

    return (
        <main className="max-w-md mx-auto">
            <TaskForm addTask={addTask} />

            <br />

            <TaskList tasks={tasks} />
        </main>
    );
}

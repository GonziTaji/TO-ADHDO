import { useState } from 'react';
import { v4 as uuidv4 } from 'uuid';

interface Task {
    id: string;
    name: string;
    completed: boolean;
}

export default function Home() {
    const [tasks, setTasks] = useState<Task[]>([]);
    const [newTaskName, setNewTaskName] = useState('');

    function addNewTask(name: string) {
        setTasks([...tasks, {
            id: uuidv4(),
            name,
            completed: false
        }]);

        setNewTaskName('');
    }

    function changeTaskName(id: string, newName: string) {
        setTasks(_tasks => _tasks.map((task) => {
            if (task.id === id) {
                task.name = newName;
            }

            return task;
        }));
    }

    function toggleTaskCompleted(id: string) {
        setTasks(_tasks => _tasks.map((task) => {
            if (task.id === id) {
                task.completed = !task.completed;
            }

            return task;
        }));
    }

    return (
        <main style={{minWidth: '1vh'}}>
            <input onKeyDown={(ev) => ev.code === 'Enter' && addNewTask(newTaskName)} type="text" value={newTaskName} onChange={(ev) => setNewTaskName(ev.currentTarget.value)} />
            <button type="button" onClick={() => addNewTask(newTaskName)}>
                Create new Task
            </button>

            <ul>
                {tasks.map((task, i) => (
                    <li key={i}>
                        <input type="checkbox" checked={task.completed} onChange={() => toggleTaskCompleted(task.id)} />
                        <input type="text" value={task.name} onChange={(ev) => changeTaskName(task.id, ev.currentTarget.value)} />
                    </li>
                ))}
            </ul>
        </main>
    );
}

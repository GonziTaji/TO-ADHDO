import TagSelect from '@/components/TagSelect';
import { useState } from 'react';
import { v4 as uuidv4 } from 'uuid';

interface Task {
    id: string;
    name: string;
    tags: string[];
    completed: boolean;
}

export default function Home() {
    const [tasks, setTasks] = useState<Task[]>([]);
    const [newTaskName, setNewTaskName] = useState('');
    const [newTaskTags, setNewTaskTags] = useState<string[]>([]);

    function addNewTask() {
        if (!newTaskName) {
            return;
        }

        setTasks([
            ...tasks,
            {
                id: uuidv4(),
                name: newTaskName,
                tags: newTaskTags,
                completed: false,
            },
        ]);

        setNewTaskName('');
        setNewTaskTags([]);
    }

    function changeTaskName(id: string, newName: string) {
        setTasks((_tasks) =>
            _tasks.map((task) => {
                if (task.id === id) {
                    task.name = newName;
                }

                return task;
            })
        );
    }

    function toggleTaskCompleted(id: string) {
        setTasks((_tasks) =>
            _tasks.map((task) => {
                if (task.id === id) {
                    task.completed = !task.completed;
                }

                return task;
            })
        );
    }

    function addTag(newTag: string) {
        if (!newTag || newTaskTags.includes(newTag)) {
            return;
        }

        const newTags = [...newTaskTags, newTag];

        setNewTaskTags(newTags);
    }

    function removeTag(tag: string) {
        const newTags = newTaskTags.filter((t) => t !== tag);

        setNewTaskTags(newTags);
    }

    return (
        <main className="max-w-md">
            <form>
                <input
                    className="border border-gray-400 px-2 py-1"
                    type="text"
                    value={newTaskName}
                    onChange={(ev) => setNewTaskName(ev.currentTarget.value)}
                    placeholder="What needs to be done?"
                />

                <div className="flex my-2">
                    {newTaskTags.map((tag, i) => (
                        <div
                            key={i}
                            className="flex gap-2 px-3 py-2 border rounded-md border-blue-400 bg-blue-200"
                        >
                            <span>{tag}</span>

                            <span
                                className="font-bold cursor-pointer"
                                onClick={() => removeTag(tag)}
                            >
                                X
                            </span>
                        </div>
                    ))}
                </div>

                <TagSelect onSelection={addTag} />

                <button
                    className="border disabled:bg-emerald-300 disabled:text-neutral-700 disabled:cursor-not-allowed border-teal-900 px-2 py-1 bg-emerald-500"
                    type="button"
                    onClick={() => addNewTask()}
                    disabled={!newTaskName.length}
                >
                    Create new Task
                </button>
            </form>

            <ul>
                {tasks.map((task, i) => (
                    <li key={i}>
                        <input
                            type="checkbox"
                            checked={task.completed}
                            onChange={() => toggleTaskCompleted(task.id)}
                        />
                        <input
                            type="text"
                            value={task.name}
                            onChange={(ev) =>
                                changeTaskName(task.id, ev.currentTarget.value)
                            }
                        />
                    </li>
                ))}
            </ul>
        </main>
    );
}

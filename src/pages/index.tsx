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
        const newTags = [...newTaskTags, newTag];

        setNewTaskTags(newTags);
    }

    function removeTag(tag: string) {
        const newTags = newTaskTags.filter((t) => t !== tag);

        setNewTaskTags(newTags);
    }

    return (
        <main style={{ minWidth: '1vh' }}>
            <form>
                {
                    // onKeyDown={(ev) => ev.code === 'Enter' && addNewTask()}>
                }
                <input
                    type="text"
                    value={newTaskName}
                    onChange={(ev) => setNewTaskName(ev.currentTarget.value)}
                />

                <div>
                    {newTaskTags.map((tag, i) => (
                        <div key={i}>
                            <span onClick={() => removeTag(tag)}>
                                <b>X</b>
                            </span>
                            {tag}
                        </div>
                    ))}
                </div>

                <TagSelect onSelection={addTag} />

                <button type="button" onClick={() => addNewTask()}>
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

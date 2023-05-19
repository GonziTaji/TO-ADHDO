'use client';
import { useState } from 'react';
import TagSelect from './TagSelect';
import { v4 as uuidv4 } from 'uuid';
import { Task } from '@/types';

interface TaskFormProps {
    tags: any[];
}

export default function TaskForm({ tags }: TaskFormProps) {
    const [newTaskName, setNewTaskName] = useState('');
    const [newTaskTags, setNewTaskTags] = useState<string[]>([]);

    function addNewTask() {
        if (!newTaskName) {
            return;
        }

        const newTag = {
            id: uuidv4(),
            name: newTaskName,
            tags: newTaskTags,
            completed: false,
        };

        setNewTaskName('');
        setNewTaskTags([]);
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
        <form className="border-2 rounded border-slate-400 bg-slate-200 p-4 pt-0">
            <h1 className="text-xl py-2">Create new Task</h1>

            <div className="flex gap-2 flex-col">
                <span>What needs to be done?</span>

                <input
                    className="border border-gray-400 px-2 py-1"
                    type="text"
                    value={newTaskName}
                    onChange={(ev) => setNewTaskName(ev.currentTarget.value)}
                    placeholder="Watch an anime"
                />
            </div>

            <span className="block pt-4">Any keywords for the task?</span>

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

            <TagSelect onSelection={addTag} tags={tags.map((t) => t.name)} />

            <button
                className="border rounded cursor-pointer disabled:bg-green-300 disabled:text-neutral-700 disabled:cursor-not-allowed border-teal-900 px-2 py-1 bg-emerald-500"
                type="button"
                onClick={() => addNewTask()}
                disabled={!newTaskName.length}
            >
                Create new Task
            </button>
        </form>
    );
}

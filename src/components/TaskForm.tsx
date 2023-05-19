'use client';
import { Suspense, useState, useTransition } from 'react';
import TagSelect from './TagSelect';
import { Tag } from '@/prismaUtils';
import { useRouter } from 'next/navigation';

interface TaskFormProps {
    tags: Tag[];
}

export default function TaskForm({ tags }: TaskFormProps) {
    const [newTaskName, setNewTaskName] = useState('');
    const [newTaskTags, setNewTaskTags] = useState<Tag[]>([]);

    const [isPending, startTransition] = useTransition();
    const [isFetching, setIsFetching] = useState(false);
    const isMutating = isFetching || isPending;

    const router = useRouter();

    async function addNewTask() {
        if (!newTaskName) {
            return;
        }

        setIsFetching(true);

        const newTask = {
            name: newTaskName,
            tags: newTaskTags.map((t) => t.id),
            user_id: 1,
        };

        await fetch('/api/tasks', {
            method: 'POST',
            body: JSON.stringify(newTask),
        });

        startTransition(() => {
            router.refresh();
        });

        setNewTaskName('');
        setNewTaskTags([]);

        setIsFetching(false);
    }

    function addTag(tagId: number) {
        if (!tagId) {
            return;
        }

        let tag = newTaskTags.find((t) => t.id === tagId);

        if (tag) {
            // tag already selected
            return;
        }

        tag = tags.find((t) => t.id === tagId);

        if (!tag) {
            // tag doesn't exist
            alert("Tag selected doesn't exist. Please try again");
            return;
        }

        const newTags = [...newTaskTags, tag];

        setNewTaskTags(newTags);
    }

    function removeTag(tagId: number) {
        const newTags = newTaskTags.filter((t) => t.id !== tagId);

        setNewTaskTags(newTags);
    }

    if (isMutating) {
        return <p>Creating task...</p>;
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
                        <span>{tag.name}</span>

                        <span
                            className="font-bold cursor-pointer"
                            onClick={() => removeTag(tag.id)}
                        >
                            X
                        </span>
                    </div>
                ))}
            </div>

            <TagSelect onSelection={addTag} tags={tags} />

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

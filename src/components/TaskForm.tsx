'use client';
import { Suspense, useState, useTransition } from 'react';
import TagSelect from './TagSelect';
import { Tag } from '@/prismaUtils';
import { useRouter } from 'next/navigation';
import TagForm from './TagForm';
import Collapsable from './Collapsable';
import TagItem from '@/app/tasks/TagItem';

interface TaskFormProps {
    tags: Tag[];
}

export default function TaskForm({ tags }: TaskFormProps) {
    const [newTaskName, setNewTaskName] = useState('');
    const [newTaskTags, setNewTaskTags] = useState<Tag[]>([]);

    const [isPending, startTransition] = useTransition();
    const [isFetching, setIsFetching] = useState(false);
    const isMutating = isFetching || isPending;

    const [showCreateTagForm, setShowCreateTagForm] = useState(false);

    const router = useRouter();

    async function addNewTask() {
        if (!newTaskName) {
            return;
        }

        if (!newTaskTags.length) {
            alert('Add least one tag for this task');
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

            <div className="flex flex-col">
                <span>What needs to be done?</span>

                <input
                    className="border border-gray-400 px-2 py-1"
                    type="text"
                    value={newTaskName}
                    onChange={(ev) => setNewTaskName(ev.currentTarget.value)}
                    placeholder="Watch an anime"
                />
            </div>

            <div className="pt-2">
                <span>Add tags for the task:</span>
                <div className="flex gap-2">
                    <div className="grow">
                        <TagSelect onSelection={addTag} tags={tags} />
                    </div>

                    <button
                        type="button"
                        title={
                            showCreateTagForm
                                ? 'Crate new tag'
                                : 'Hide tag creation'
                        }
                        className="text-mono cursor-pointer font-bold text-xl border border-sky-900 w-8 bg-sky-200"
                        onClick={() => setShowCreateTagForm(!showCreateTagForm)}
                    >
                        {showCreateTagForm ? '-' : '+'}
                    </button>
                </div>
            </div>

            <Collapsable collapsed={!showCreateTagForm} className="mt-1">
                <div className="border border-gray-400 p-2 shadow">
                    <span>Create new Tag:</span>
                    <TagForm
                        userId={1}
                        onSubmit={() => setShowCreateTagForm(false)}
                    />
                </div>
            </Collapsable>

            <div className="flex flex-wrap gap-1 my-2">
                {newTaskTags.map((tag) => (
                    <TagItem
                        key={tag.id}
                        tag={tag}
                        onDelete={() => removeTag(tag.id)}
                    />
                ))}

                {!newTaskTags.length && (
                    <div className="text-gray-400 py-1 border">
                        <i>No tags selected</i>
                    </div>
                )}
            </div>

            <div className="flex justify-end pt-4">
                <button
                    className="inline-flex whitespace-nowrap items-center gap-3 border rounded cursor-pointer disabled:cursor-not-allowed disabled:bg-green-100 disabled:text-gray-400 border-gray-400 bg-green-400 px-2 py-1"
                    type="button"
                    onClick={() => addNewTask()}
                    disabled={!newTaskName.length}
                >
                    Create new Task
                </button>
            </div>
        </form>
    );
}

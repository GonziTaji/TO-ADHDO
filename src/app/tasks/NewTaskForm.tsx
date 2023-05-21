'use client';

import { Tags } from '@prisma/client';
import { useRouter } from 'next/navigation';
import { ChangeEvent, useState, useTransition } from 'react';
import TagSelect from './TagSelect';
import TagItem from './TagItem';
import { Spinner } from '@/components/Spinner';

interface NewTaskFormProps {
    tags: Tags[];
}

export default function NewTaskForm({ tags }: NewTaskFormProps) {
    const [taskName, setTaskName] = useState('');
    const [taskTags, setTaskTags] = useState<Tags[]>([]);

    const [isCreating, setIsCreating] = useState(false);
    const [isPending, startTransition] = useTransition();

    const isMutating = isCreating || isPending;

    const router = useRouter();

    function inputOnChange(ev: ChangeEvent<HTMLInputElement>) {
        setTaskName(ev.currentTarget.value);
    }

    async function createTask() {
        if (!taskName) {
            return;
        }

        if (!taskTags.length) {
            alert('Add least one tag for this task');
            return;
        }

        setIsCreating(true);

        const newTask = {
            name: taskName,
            tags: taskTags.map((t) => t.id),
            user_id: 1,
        };

        await fetch('/api/tasks', {
            method: 'POST',
            body: JSON.stringify(newTask),
        });

        startTransition(() => {
            setIsCreating(false);
            setTaskName('');
            setTaskTags([]);

            router.refresh();
        });
    }

    function addTag(tagId: number) {
        if (!tagId) {
            return;
        }

        let tag = taskTags.find((t) => t.id === tagId);

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

        const newTags = [...taskTags, tag];

        setTaskTags(newTags);
    }

    function removeTag(tagId: number) {
        const newTags = taskTags.filter((t) => t.id !== tagId);

        setTaskTags(newTags);
    }

    return (
        <div>
            <label htmlFor="task-name" className="mb-2 block">
                Create a new Task
            </label>

            <div className="flex gap-2 justify-between">
                <input
                    id="task-name"
                    className="grow min-w-0 ps-2"
                    value={taskName}
                    onChange={inputOnChange}
                    disabled={isMutating}
                    placeholder="What needs to be done?"
                />

                <button
                    className="inline-flex whitespace-nowrap items-center gap-3 border rounded cursor-pointer disabled:cursor-not-allowed disabled:bg-green-100 disabled:text-gray-400 border-gray-400 bg-green-400 px-2 py-1"
                    type="button"
                    onClick={() => createTask()}
                    disabled={!taskName.length || isMutating}
                >
                    Save Task {isMutating && <Spinner />}
                </button>
            </div>

            <ul className="flex flex-wrap gap-1 my-2 items-center">
                <div className="max-w-[7rem]">
                    <TagSelect tags={tags} onSelection={addTag} />
                </div>

                {!taskTags.length && (
                    <li className="text-xs border rounded-full border-orange-100 bg-amber-100 text-gray-600 px-2 m-1">
                        No tags
                    </li>
                )}

                {taskTags.map((tag) => (
                    <TagItem key={tag.id} tag={tag} onDelete={removeTag} />
                ))}
            </ul>
        </div>
    );
}

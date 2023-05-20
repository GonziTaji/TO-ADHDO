'use client';

import { TaskWithTags } from '@/prismaUtils';
import { useRouter } from 'next/navigation';
import { useState, useTransition } from 'react';
import TagItem from './TagItem';
import TagSelect from '@/components/TagSelect';
import { Tags } from '@prisma/client';
import { Spinner } from '@/components/Spinner';

interface TaskItemProps {
    task: TaskWithTags;
    tags: Tags[];
}

export default function TaskItemList({ task, tags }: TaskItemProps) {
    const [isDeleting, setIsDeleting] = useState(false);
    const [isAdding, setIsAdding] = useState(false);
    const [isDeletionPending, startDeletionTransition] = useTransition();
    const [isAdditionPending, startAdditionTransition] = useTransition();

    const isMutatingTag = isAdding || isAdditionPending;

    const [isEditMode, setIsEditMode] = useState(false);

    const router = useRouter();

    async function deleteTask() {
        const confirmed = confirm(
            `Do you really want to delete the task "${task.name}"?`
        );

        if (!confirmed) {
            return;
        }

        setIsDeleting(true);

        const response = await fetch('/api/tasks/' + task.id, {
            method: 'DELETE',
        });

        const body = await response.json();

        console.log(body);

        startDeletionTransition(() => {
            router.refresh();
        });

        setIsDeleting(false);
    }

    async function addTag(tagId: number) {
        if (task.Tags.find(({ id }) => id === tagId)) {
            return;
        }

        setIsAdding(true);

        const url = `/api/tags_of_tasks?task_id=${task.id}&tag_id=${tagId}`;

        const response = await fetch(url, {
            method: 'POST',
            body: JSON.stringify({
                task_id: task.id,
                tag_id: tagId,
            }),
        });

        response.json().then(console.log);

        startAdditionTransition(() => {
            setIsAdding(false);
            router.refresh();
        });
    }

    function toggleEditMode() {
        setIsEditMode(!isEditMode);
    }

    if (isDeleting || isDeletionPending) {
        return <p>Deleting task...</p>;
    }

    return (
        <div className="px-4 py-1 border-b border-slate-500">
            <div className="flex justify-between">
                <span>{task.name}</span>

                <div className="flex gap-2 items-center">
                    <button
                        type="button"
                        className="h-7 cursor-pointer bg-indigo-200 rounded px-3 border border-slate-400"
                        onClick={toggleEditMode}
                    >
                        Edit
                    </button>

                    <button
                        type="button"
                        className="h-7 cursor-pointer text-rose-900 bg-red-200 rounded px-3 border border-slate-400"
                        onClick={deleteTask}
                    >
                        Delete
                    </button>
                </div>
            </div>

            <ul className="flex flex-wrap gap-1 my-2">
                <li>Tags: </li>

                {!task.Tags.length && (
                    <li className="text-xs border rounded-full border-orange-100 bg-amber-100 text-gray-600 px-2 m-1">
                        No tags
                    </li>
                )}

                {task.Tags.map((tag) => (
                    <TagItem key={tag.id} tag={tag} taskId={task.id} />
                ))}

                <div className="inline-flex gap-1 items-center">
                    <div className="max-w-[7rem]">
                        <TagSelect
                            onSelection={addTag}
                            tags={tags}
                            disabled={isMutatingTag}
                        />
                    </div>
                    {isMutatingTag && <Spinner />}
                </div>
            </ul>
        </div>
    );
}

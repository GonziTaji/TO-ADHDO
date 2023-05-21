'use client';

import { TaskWithTags } from '@/prismaUtils';
import { useRouter } from 'next/navigation';
import { useState, useTransition } from 'react';
import TagItem from './TagItem';
import TagSelect from './TagSelect';
import { Tags } from '@prisma/client';
import TaskNameInput from './TaskNameInput';

interface TaskItemProps {
    task: TaskWithTags;
    tags: Tags[];
}

export default function TaskItemList({ task, tags }: TaskItemProps) {
    const [isDeleting, setIsDeleting] = useState(false);
    const [isDeletionPending, startDeletionTransition] = useTransition();

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

    if (isDeleting || isDeletionPending) {
        return <p>Deleting task...</p>;
    }

    return (
        <div>
            <div className="flex gap-2 justify-between">
                <TaskNameInput task={task} />

                <button
                    type="button"
                    className="h-7 cursor-pointer text-rose-900 bg-red-200 rounded px-3 border border-slate-400"
                    onClick={deleteTask}
                >
                    Delete
                </button>
            </div>

            <ul className="flex flex-wrap gap-1 my-2">
                <li>Tags: </li>

                <li className="max-w-[7rem]">
                    <TagSelect tags={tags} task={task} />
                </li>

                {!task.Tags.length && (
                    <li className="text-xs border rounded-full border-orange-100 bg-amber-100 text-gray-600 px-2 m-1">
                        No tags
                    </li>
                )}

                {task.Tags.map((tag) => (
                    <TagItem key={tag.id} tag={tag} taskId={task.id} />
                ))}
            </ul>
        </div>
    );
}

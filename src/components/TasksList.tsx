'use client';

import { TaskWithTags } from '@/prismaUtils';
import { useRouter } from 'next/navigation';
import { useState, useTransition } from 'react';

export function TaskList({ tasks }: { tasks: TaskWithTags[] }) {
    const [isPending, startTransition] = useTransition();
    const [isFetching, setIsFetching] = useState(false);
    const isMutating = isFetching || isPending;

    const router = useRouter();

    async function deleteTask(task: TaskWithTags) {
        const confirmed = confirm(
            `Do you really want to delete the task ${task.name}?`
        );

        if (!confirmed) {
            return;
        }

        setIsFetching(true);

        const response = await fetch('/api/tasks/' + task.id, {
            method: 'DELETE',
        });

        const body = response.json();

        console.log(body);

        startTransition(() => {
            router.refresh();
        });

        setIsFetching(false);
    }

    if (isMutating) {
        return <p>Loading...</p>;
    }

    return (
        <ul className="flex flex-col pt-2">
            {tasks.map((task, i) => (
                <li key={i} className="border-b border-black">
                    <div className="flex justify-between">
                        <span className="px-2 py-1">{task.name}</span>
                        <button
                            type="button"
                            className="cursor-pointer"
                            onClick={() => deleteTask(task)}
                        >
                            Delete
                        </button>
                    </div>

                    <ul className="ms-4 flex">
                        {task.Tags.map((tag) => (
                            <li
                                key={tag.id}
                                className="border border-rose-400 rounded bg-amber-300 px-2 m-1"
                            >
                                {tag.name}
                            </li>
                        ))}

                        {!task.Tags.length && (
                            <li className="border border-rose-300 rounded bg-amber-200 text-gray-600 px-2 m-1">
                                No tags
                            </li>
                        )}
                    </ul>
                </li>
            ))}
        </ul>
    );
}

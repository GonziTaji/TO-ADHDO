'use client';

import { TaskWithTags } from '@/prismaUtils';
import { useRouter } from 'next/navigation';
import { useState, useTransition } from 'react';
import TaskItemList from './TaskItemList';
import { Tags } from '@prisma/client';

interface TaskListProps {
    tasks: TaskWithTags[];
    tags: Tags[];
}

export function TaskList({ tasks, tags }: TaskListProps) {
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
        <div className="flex flex-col pt-2">
            {tasks.map((task, i) => (
                <TaskItemList task={task} key={i} tags={tags} />
            ))}
        </div>
    );
}

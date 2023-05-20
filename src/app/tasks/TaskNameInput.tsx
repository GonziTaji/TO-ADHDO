import { Spinner } from '@/components/Spinner';
import { TaskWithTags } from '@/prismaUtils';
import { useRouter } from 'next/navigation';
import { ChangeEvent, useState, useTransition } from 'react';

interface TaskNameInputProps {
    task: TaskWithTags;
}

export default function TaskNameInput({ task }: TaskNameInputProps) {
    const [taskName, setTaskName] = useState(task?.name || '');

    const [isUpdating, setIsUpdating] = useState(false);
    const [isPending, starTransition] = useTransition();

    const isMutating = isUpdating || isPending;

    const router = useRouter();

    function inputOnChange(ev: ChangeEvent<HTMLInputElement>) {
        console.log('onChange');
        setTaskName(ev.currentTarget.value);
    }

    async function inputOnBlur() {
        if (taskName.trim() === task.name) {
            return;
        }

        setIsUpdating(true);

        await fetch('/api/tasks/' + task.id, {
            method: 'PUT',
            body: JSON.stringify({
                name: taskName,
            }),
        });

        starTransition(() => {
            setIsUpdating(false);
            router.refresh();
        });
    }

    return (
        <div className="contents">
            <input
                className="min-w-0"
                value={taskName}
                onChange={inputOnChange}
                onBlur={inputOnBlur}
                disabled={isMutating}
            />
        </div>
    );
}

import { Spinner } from '@/components/Spinner';
import { Tags } from '@prisma/client';
import { useRouter } from 'next/navigation';
import { useState, useTransition } from 'react';

interface TagItemProps {
    taskId?: number;
    tag: Tags;
    onDelete?: (tagId: number) => void;
}

export default function TagItem({ taskId, tag, onDelete }: TagItemProps) {
    const [isPending, startTransition] = useTransition();
    const [isFetching, setIsFetching] = useState(false);
    const isMutating = isFetching || isPending;

    const router = useRouter();

    async function deleteTag() {
        if (taskId) {
            setIsFetching(true);

            const url = `/api/tags_of_tasks?task_id=${taskId}&tag_id=${tag.id}`;

            const response = await fetch(url, { method: 'DELETE' });

            response.json().then(console.log);

            startTransition(() => {
                setIsFetching(false);
                router.refresh();
            });
        }

        if (onDelete) {
            onDelete(tag.id);
        }
    }

    return (
        <div
            className={
                'text-xs max-w-[7rem] flex items-center justify-around gap-2 px-2 py-1 border rounded-full ' +
                (isMutating
                    ? 'border-orange-100 bg-amber-100 text-gray-600'
                    : 'border-orange-200 bg-amber-200')
            }
        >
            <span
                title={tag.name}
                className="cursor-default grow text-center overflow-hidden overflow-ellipsis whitespace-nowrap"
            >
                {tag.name}
            </span>

            {isMutating ? (
                <Spinner />
            ) : (
                <span
                    className="font-bold font-mono cursor-pointer"
                    onClick={deleteTag}
                >
                    x
                </span>
            )}
        </div>
    );
}

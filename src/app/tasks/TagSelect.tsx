import { Spinner } from '@/components/Spinner';
import { Tag, TaskWithTags } from '@/prismaUtils';
import { useRouter } from 'next/navigation';
import { ChangeEvent, useState, useTransition } from 'react';

interface TagSelectProps {
    tags: Tag[];
    task?: TaskWithTags;
    onSelection?: (tagId: number) => void;
}

export default function TagSelect({ tags, task, onSelection }: TagSelectProps) {
    const [selection, setSelection] = useState('');

    const [isAdding, setIsAdding] = useState(false);
    const [isPending, starTransition] = useTransition();

    const isMutating = isAdding || isPending;

    const router = useRouter();

    async function selectOnChange(ev: ChangeEvent<HTMLSelectElement>) {
        setSelection(ev.currentTarget.value);

        const tagId = parseInt(ev.currentTarget.value);

        if (isNaN(tagId)) {
            return;
        }

        if (task) {
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

            starTransition(() => {
                setIsAdding(false);
                router.refresh();
            });
        }

        if (onSelection) {
            onSelection(tagId);
        }

        setSelection('');
    }

    return (
        <div className="flex gap-2 items-center text-sm">
            <select
                id="tag-select"
                className="grow border border-gray-400 px-2 py-1 min-w-0 rounded-full"
                onChange={selectOnChange}
                disabled={isMutating}
                value={selection}
            >
                <option value="">Add tag</option>
                {tags.map((tag, i) => (
                    <option key={i} value={tag.id}>
                        {tag.name}
                    </option>
                ))}
            </select>
            {isMutating && <Spinner />}
        </div>
    );
}

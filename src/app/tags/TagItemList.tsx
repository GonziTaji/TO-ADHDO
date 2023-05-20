'use client';

import TagForm from '@/components/TagForm';
import { TagWithTaskCount } from '@/prismaUtils';
import { useRouter } from 'next/navigation';
import { useState, useTransition } from 'react';

interface TagItemListProps {
    tag: TagWithTaskCount;
}

export default function TagItemList({ tag }: TagItemListProps) {
    const [isDeleting, setIsDeleting] = useState(false);
    const [isPending, startTransition] = useTransition();
    const isMutating = isDeleting || isPending;

    const [isEditMode, setIsEditMode] = useState(false);

    const router = useRouter();

    async function deleteTag() {
        const confirmed = confirm(
            `Do you really want to delete the tag "${tag.name}"?`
        );

        if (!confirmed) {
            return;
        }

        setIsDeleting(true);

        const response = await fetch('/api/tags/' + tag.id, {
            method: 'DELETE',
        });

        const body = await response.json();

        console.log(body);

        startTransition(() => {
            router.refresh();
        });

        setIsDeleting(false);
    }

    function toggleEditMode() {
        setIsEditMode(!isEditMode);
    }

    if (isMutating) {
        return <p>Deleting tag...</p>;
    }

    const componentBody = isEditMode ? (
        <div className="gap-2 px-4 py-1">
            <TagForm
                userId={1}
                tagId={tag.id}
                tagName={tag.name}
                onSubmit={toggleEditMode}
            />

            <button
                type="button"
                className="mt-2 cursor-pointer text-rose-900 bg-red-200 rounded px-3 border border-slate-400"
                onClick={toggleEditMode}
            >
                Cancel
            </button>
        </div>
    ) : (
        <div className="flex justify-between px-4 py-1">
            <span className="whitespace-nowrap overflow-ellipsis overflow-hidden">
                <span title={tag.name}>{tag.name}</span>

                <small className="block">
                    Present in {tag.task_count} task
                    {tag.task_count !== 1 && 's'}
                </small>
            </span>

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
                    onClick={deleteTag}
                >
                    Delete
                </button>
            </div>
        </div>
    );

    return <div className="border-b border-slate-500">{componentBody}</div>;
}

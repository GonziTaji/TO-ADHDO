'use client';

import TagForm from '@/components/TagForm';
import { TagWithTaskCount } from '@/prismaUtils';
import { useRouter } from 'next/navigation';
import { useState, useTransition } from 'react';

interface TagItemProps {
    tag: TagWithTaskCount;
}

export default function TagItem({ tag }: TagItemProps) {
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
        <div className="flex px-2">
            <TagForm
                userId={1}
                tagId={tag.id}
                tagName={tag.name}
                onSubmit={toggleEditMode}
            />

            <button
                type="button"
                className="cursor-pointer text-rose-900 bg-red-200 rounded px-3 border border-slate-400"
                onClick={toggleEditMode}
            >
                Cancel
            </button>
        </div>
    ) : (
        <div
            className="grid px-4 py-1 "
            style={{ gridTemplateColumns: '1fr auto' }}
        >
            <span>
                {tag.name}

                <small className="block">
                    Present in {tag.task_count} task
                    {tag.task_count !== 1 && 's'}
                </small>
            </span>

            <div className="flex gap-2">
                <button
                    type="button"
                    className="cursor-pointer bg-indigo-200 rounded px-3 border border-slate-400"
                    onClick={toggleEditMode}
                >
                    Edit
                </button>

                <button
                    type="button"
                    className="cursor-pointer text-rose-900 bg-red-200 rounded px-3 border border-slate-400"
                    onClick={deleteTag}
                >
                    Delete
                </button>
            </div>
        </div>
    );

    return <div className="border-b border-slate-500">{componentBody}</div>;
}

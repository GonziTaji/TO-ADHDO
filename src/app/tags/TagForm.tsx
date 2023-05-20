'use client';

import { useRouter } from 'next/navigation';
import { KeyboardEvent, useState, useTransition } from 'react';
import { Spinner } from '../../components/Spinner';

interface TagFormProps {
    userId: number;
    tagId?: number;
    tagName?: string;
    onSubmit?: () => void;
}

export default function TagForm({
    userId,
    tagId,
    tagName,
    onSubmit,
}: TagFormProps) {
    const [isCreating, setIsCreating] = useState(false);
    const [isPending, startTransition] = useTransition();
    const isMutating = isCreating || isPending;

    const [newTagName, setNewTagName] = useState(tagName || '');

    const router = useRouter();

    async function saveChanges() {
        if (!newTagName) {
            return;
        }

        setIsCreating(true);

        const newTag = {
            name: newTagName,
            user_id: userId,
        };

        let url = '/api/tags/';
        let method = 'POST';

        if (tagId) {
            url += tagId;
            method = 'PUT';
        }

        const response = await fetch(url, {
            method,
            body: JSON.stringify(newTag),
        });

        response.json().then(console.log);

        startTransition(() => {
            setNewTagName('');
            setIsCreating(false);

            if (onSubmit && typeof onSubmit === 'function') {
                onSubmit();
            }

            router.refresh();
        });
    }

    function onKeyDownInput(ev: KeyboardEvent) {
        if (ev.code === 'Enter') {
            saveChanges();
        }
    }

    return (
        <div className="flex gap-2 items-center">
            <input
                id="tag-name"
                className="grow min-w-0 border border-gray-400 px-2 py-1"
                type="text"
                value={newTagName}
                onInput={(ev) => setNewTagName(ev.currentTarget.value)}
                onKeyDown={onKeyDownInput}
                placeholder="food/cleaning/kitchen/etc"
                disabled={isMutating}
            />

            <button
                className="inline-flex whitespace-nowrap items-center gap-3 border rounded cursor-pointer disabled:cursor-not-allowed disabled:bg-green-100 disabled:text-gray-400 border-gray-400 bg-green-400 px-2 py-1"
                type="button"
                onClick={saveChanges}
                disabled={isMutating || !newTagName.trim().length}
            >
                {tagId ? 'Update' : 'Add'} tag {isMutating && <Spinner />}
            </button>
        </div>
    );
}

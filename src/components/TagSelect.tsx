import { Tag } from '@/prismaUtils';
import { useRouter } from 'next/navigation';
import { ChangeEvent, KeyboardEvent, useState, useTransition } from 'react';

interface TagSelectProps {
    onSelection: (tagId: number) => void;
    tags: Tag[];
}

export default function TagSelect({ onSelection, tags }: TagSelectProps) {
    const [newTag, setNewTag] = useState('');

    const [isPending, startTransition] = useTransition();
    const [isFetching, setIsFetching] = useState(false);
    const isMutating = isFetching || isPending;

    const router = useRouter();

    async function createTag() {
        if (!newTag) {
            return;
        }

        setIsFetching(true);

        const tagFound = tags.find((t) => t.name === newTag);

        if (!tagFound) {
            await fetch('/api/tags', {
                method: 'POST',
                body: JSON.stringify({
                    name: newTag,
                    user_id: 1,
                }),
            });

            startTransition(() => {
                router.refresh();
            });
        } else {
            alert(`Tag "${newTag}" already exist!`);
        }

        setIsFetching(false);

        setNewTag('');
    }

    function onKeyDownInput(ev: KeyboardEvent) {
        if (ev.code === 'Enter') {
            createTag();
        }
    }

    function selectOnChange(ev: ChangeEvent<HTMLSelectElement>) {
        const value = parseInt(ev.currentTarget.value);

        if (!isNaN(value)) {
            onSelection(value);
        }
    }

    if (isMutating) {
        return <p>Creating tag...</p>;
    }

    return (
        <div className="flex flex-col gap-3">
            <div className="flex gap-2 items-center">
                <select
                    id="tag-select"
                    className="grow border border-gray-400 px-2 py-1"
                    onChange={selectOnChange}
                >
                    <option value="">-- Select an option</option>
                    {tags.map((tag, i) => (
                        <option key={i} value={tag.id}>
                            {tag.name}
                        </option>
                    ))}
                </select>
            </div>

            <div className="flex gap-2 items-center whitespace-nowrap flex-wrap justify-end">
                <span>Create a new one:</span>

                <input
                    id="tag-name"
                    className="grow border border-gray-400 px-2 py-1"
                    type="text"
                    value={newTag}
                    onInput={(ev) => setNewTag(ev.currentTarget.value)}
                    onKeyDown={onKeyDownInput}
                    placeholder="food/cleaning/kitchen/etc"
                />

                <button
                    className="border rounded cursor-pointer border-gray-400 bg-green-400 px-2 py-1"
                    type="button"
                    onClick={createTag}
                    disabled={!newTag.trim().length}
                >
                    Add new Tag
                </button>
            </div>
        </div>
    );
}

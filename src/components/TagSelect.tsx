import { KeyboardEvent, useState } from 'react';

export interface SelectOption {
    label: string;
    id: string;
    hovered?: boolean;
}

const __tags: string[] = [];

interface TagSelectProps {
    onSelection: (selectedId: string) => void;
    tags: string[];
}

export default function TagSelect({ onSelection, tags }: TagSelectProps) {
    const [newTag, setNewTag] = useState('');

    function setKnownTags(p: any) {}

    async function createTag() {
        if (!newTag) {
            return;
        }

        if (!tags.includes(newTag)) {
            setKnownTags([...tags, newTag]);
        } else {
            alert(`Tag "${newTag}" already exist!`);
        }

        onSelection(newTag);
        setNewTag('');
    }

    function onKeyDownInput(ev: KeyboardEvent) {
        if (ev.code === 'Enter') {
            createTag();
        }
    }

    return (
        <div className="flex flex-col gap-3">
            <div className="flex gap-2 items-center">
                <select
                    id="tag-select"
                    className="grow border border-gray-400 px-2 py-1"
                    onChange={(ev) => onSelection(ev.currentTarget.value)}
                >
                    <option value="">-- </option>
                    {tags.map((tag, i) => (
                        <option key={i} value={tag}>
                            {tag}
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

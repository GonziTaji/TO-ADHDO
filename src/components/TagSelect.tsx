import { KeyboardEvent, useState } from 'react';

export interface SelectOption {
    label: string;
    id: string;
    hovered?: boolean;
}

const __tags: string[] = [];

interface TagSelectProps {
    onSelection: (selectedId: string) => void;
}

export default function TagSelect({ onSelection }: TagSelectProps) {
    const [knownTags, setKnownTags] = useState<string[]>(__tags);
    const [newTag, setNewTag] = useState('');

    const [createMode, setCreateMode] = useState(false);

    async function createTag() {
        if (!newTag) {
            return;
        }

        if (!knownTags.includes(newTag)) {
            setKnownTags([...knownTags, newTag]);
        } else {
            alert(`Tag "${newTag}" already exist!`);
        }

        onSelection(newTag);
        setNewTag('');
        setCreateMode(false);
    }

    function toggleCreateMode() {
        setCreateMode(!createMode);
    }

    return (
        <div className="">
            {createMode ? (
                <>
                    <div className="flex flex-col">
                        <label htmlFor="tag-name">Create a new Task</label>
                        <input
                            id="tag-name"
                            className="border border-gray-400 px-2 py-1"
                            type="text"
                            value={newTag}
                            onInput={(ev) => {
                                setNewTag(ev.currentTarget.value);
                            }}
                            onKeyDown={(ev) =>
                                ev.code === 'Enter' && createTag()
                            }
                            placeholder="food/cleaning/kitchen/etc"
                        />
                    </div>

                    <div className="flex gap-2">
                        <button
                            className="border border-gray-400 px-2 py-1"
                            type="button"
                            style={{ height: '28px', padding: '5px 3px' }}
                            onClick={createTag}
                            disabled={!newTag.trim().length}
                        >
                            Add new Tag
                        </button>

                        <button
                            className="border border-gray-400 px-2 py-1"
                            type="button"
                            onClick={toggleCreateMode}
                        >
                            Cancel
                        </button>
                    </div>
                </>
            ) : (
                <div className="flex flex-col gap-3">
                    <label htmlFor="tag-select">
                        Select tags for this Task
                    </label>
                    <select
                        id="tag-select"
                        onChange={(ev) => onSelection(ev.currentTarget.value)}
                    >
                        <option value="">Select a tag</option>
                        {knownTags.map((tag, i) => (
                            <option key={i} value={tag}>
                                {tag}
                            </option>
                        ))}
                    </select>

                    <button
                        className="border border-gray-400 px-2 py-1"
                        type="button"
                        onClick={toggleCreateMode}
                    >
                        Create new Tag
                    </button>
                </div>
            )}
        </div>
    );
}

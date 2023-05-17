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
    const [searchTerm, setSearchTerm] = useState('');
    const [hoveredOptionId, setHoveredOptionId] = useState('');
    const [showResults, setShowResults] = useState(false);
    const [isMouseDownParent, setIsMouseDownParent] = useState(false);

    const tagOptions = knownTags
        .filter((tag) =>
            searchTerm ? RegExp(searchTerm, 'i').test(tag) : true
        )
        .map((tag) => ({
            label: tag,
            id: tag,
            hovered: false,
        }));

    async function createTag(newTag: string) {
        if (!newTag) {
            alert('tag cannot be empty');
            return;
        }

        setKnownTags([...knownTags, newTag]);
        onSelection(newTag);
        setSearchTerm('');
        setShowResults(false);
    }

    const showCreateButton = searchTerm.trim() && !tagOptions.length;

    function handleOnKeyDown(ev: KeyboardEvent) {
        if (ev.code === 'Enter') {
            if (showCreateButton) {
                createTag(searchTerm);
            } else {
                if (hoveredOptionId) {
                    selectTag(hoveredOptionId);
                }
            }
        }
    }

    function selectTag(tagId: string) {
        setSearchTerm('');
        setShowResults(false);
        onSelection(tagId);
    }

    return (
        <div
            onFocus={() => setShowResults(true)}
            // Mouse down/up to check if an option was clicked (click down inside element) or the click was outside
            // this element. If that's not handled, the click on the options is not registered, since those are
            // hidden on this element's blur
            onMouseDown={() => setIsMouseDownParent(true)}
            onBlur={() => !isMouseDownParent && setShowResults(false)}
            onMouseUp={() => setIsMouseDownParent(false)}
            onKeyDown={handleOnKeyDown}
        >
            <input
                className="border border-gray-400 px-2 py-1"
                type="text"
                value={searchTerm}
                onInput={(ev) => {
                    setSearchTerm(ev.currentTarget.value);
                }}
                placeholder="Search or create a keyword"
            />

            <div
                className="min-w-[240px] max-h-60 absolute bg-white border border-gray-400 z-10 overflow-y-scroll whitespace-nowrap"
                style={{ display: showResults ? 'block' : 'none' }}
            >
                {tagOptions.map((o_result) => {
                    const style: any = {
                        height: '28px',
                        padding: '5px 3px',
                        cursor: 'pointer',
                    };

                    if (hoveredOptionId === o_result.id) {
                        style.color = 'white';
                        style.backgroundColor = '#607799';
                    }

                    return (
                        <div
                            key={o_result.id}
                            style={style}
                            onMouseEnter={() => setHoveredOptionId(o_result.id)}
                            onClick={() => selectTag(o_result.id)}
                        >
                            {o_result.label}
                        </div>
                    );
                })}
                {showCreateButton && (
                    <button
                        type="button"
                        style={{ height: '28px', padding: '5px 3px' }}
                        onClick={() => createTag(searchTerm)}
                    >
                        Add &quot;{searchTerm}&quot;
                    </button>
                )}
            </div>
        </div>
    );
}

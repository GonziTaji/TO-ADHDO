import { Tag } from '@/prismaUtils';
import { ChangeEvent } from 'react';

interface TagSelectProps {
    onSelection: (tagId: number) => void;
    tags: Tag[];
    disabled?: boolean;
}

export default function TagSelect({
    onSelection,
    tags,
    disabled,
}: TagSelectProps) {
    function selectOnChange(ev: ChangeEvent<HTMLSelectElement>) {
        const value = parseInt(ev.currentTarget.value);

        if (!isNaN(value)) {
            onSelection(value);
        }
    }

    return (
        <div className="flex gap-2 items-center">
            <select
                id="tag-select"
                className="grow border border-gray-400 px-2 py-1 min-w-0"
                onChange={selectOnChange}
                disabled={disabled}
            >
                <option value="">-- Select an option</option>
                {tags.map((tag, i) => (
                    <option key={i} value={tag.id}>
                        {tag.name}
                    </option>
                ))}
            </select>
        </div>
    );
}

export default function TagForm() {
    const [newTag, setNewTag] = useState('');

    async function createTag() {
        if (!newTag) {
            return;
        }

        if (!knownTags.includes(newTag)) {
            setKnownTags([...knownTags, newTag]);
        } else {
            alert(`Tag "${newTag}" already exist!`);
        }

        setNewTag('');
    }

    function onKeyDownInput(ev: KeyboardEvent) {
        if (ev.code === 'Enter') {
            createTag();
        }
    }

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
    </div>;
}

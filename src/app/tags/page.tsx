import { TagWithTaskCount, getTagsWithTaskCountOfUser } from '@/prismaUtils';
import TagItem from './TagItem';
import TagForm from '@/components/TagForm';

async function getTags() {
    try {
        const tags = await getTagsWithTaskCountOfUser(1);

        return { tags };
    } catch (error) {
        return { error, tags: [] as TagWithTaskCount[] };
    }
}

export default async function Page() {
    const { tags } = await getTags();

    return (
        <div>
            <h1 className="text-2xl">Tags</h1>
            <section className="border border-rouned border-slate-400 p-4 my-2">
                <h2>Create new tag</h2>
                <TagForm userId={1} />
            </section>

            <section>
                {tags.map((tag) => (
                    <TagItem key={tag.id} tag={tag} />
                ))}
            </section>
        </div>
    );
}

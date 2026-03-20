import { Tag } from '@/shared/components'
import type { Tag as TagType } from '@/api'
import styles from './ArticleTags.module.css'

interface ArticleTagsProps {
    tags: TagType[]
    onSelect?: (id: string) => void
    selectedIds?: string[]
}

export function ArticleTags({ tags, onSelect, selectedIds = [] }: ArticleTagsProps) {
    if (!tags.length) return null

    return (
        <ul className={styles.list}>
            {tags.map((tag) => (
                <li key={tag.Id}>
                    <Tag
                        name={tag.Name}
                        selected={selectedIds.includes(tag.Id)}
                        onClick={onSelect ? () => onSelect(tag.Id) : undefined}
                    />
                </li>
            ))}
        </ul>
    )
}

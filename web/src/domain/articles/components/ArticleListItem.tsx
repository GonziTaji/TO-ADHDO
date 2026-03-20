import { type ArticleListItem as ArticleListItemType } from '@/api'
import { Button } from '@/shared/components'
import styles from './ArticleListItem.module.css'

interface ArticleListItemProps {
    article: ArticleListItemType
    onView?: (id: string) => void
    onEdit?: (id: string) => void
    onDelete?: (id: string) => void
}

export function ArticleListItemComponent({ article, onView, onEdit, onDelete }: ArticleListItemProps) {
    const firstPrice = article.Prices?.[0]?.Price ?? article.Price


    return (
        <li className={styles.row} data-component="articles-list-item">
            <span className={styles.name}>{article.Name}</span>

            <div className={styles.tagsContainer}>
                <ul className={styles.tagsList}>
                    {article.Tags.map((tag) => (
                        <li key={tag.Id}>
                            <span>{tag.Name}</span>
                        </li>
                    ))}
                </ul>
                <div>
                    <span>Price: </span>
                    <span>{firstPrice}</span>
                </div>
            </div>

            <div className={styles.buttons}>
                {onView && <Button onClick={() => onView(article.Id)}>View</Button>}
                {onEdit && <Button variant="secondary" onClick={() => onEdit(article.Id)}>Edit</Button>}
                {onDelete && <Button variant="danger" onClick={() => onDelete(article.Id)}>Delete</Button>}
            </div>
        </li>
    )
}

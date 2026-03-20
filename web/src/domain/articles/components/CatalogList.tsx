import { type ArticleListItem as ArticleListItemType } from '@/api'
import { PriceDisplay } from '@/shared/components'
import styles from './CatalogList.module.css'

interface CatalogListProps {
    articles: ArticleListItemType[]
}

export function CatalogList({ articles }: CatalogListProps) {
    return (
        <ul className={styles.list}>
            {articles.map((article) => (
                <li key={article.Id} className={styles.item}>
                    <a href={`/articles/${article.Id}.html`} className={styles.imageWrapper}>
                        <img src={article.ThumbnailUrl} alt={`${article.Name}'s image`} />
                    </a>

                    <div className={styles.info}>
                        <div className={styles.badges}>
                            {article.AvailableForTrade && (
                                <span className={styles.badge}>Permuta</span>
                            )}
                            {article.Condition && (
                                <span className={styles.badge}>{article.Condition.Label}</span>
                            )}
                        </div>

                        <h3 className={styles.name}>
                            <a href={`/articles/${article.Id}.html`}>{article.Name}</a>
                        </h3>

                        <ArticleTagList tags={article.Tags} />

                        <PriceDisplay price={article.Price} />
                    </div>
                </li>
            ))}
        </ul>
    )
}

function ArticleTagList({ tags }: { tags: ArticleListItemType['Tags'] }) {
    if (!tags.length) return null

    return (
        <ul className={styles.tags}>
            {tags.map((tag) => (
                <li key={tag.Id} className={styles.tag}>
                    {tag.Name}
                </li>
            ))}
        </ul>
    )
}

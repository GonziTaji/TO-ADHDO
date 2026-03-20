import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { Layout, Button, PriceDisplay } from '@/shared/components'
import { ArticleTags } from '../components'
import { articlesApi, type Article } from '@/api'
import { useCart } from '@/shared/context'
import { useApi } from '@/shared/hooks'
import { useState } from 'react'
import styles from './view.module.css'

async function fetchArticle() {
    // Route: /catalog/:article_id — the last path segment is the id
    const path = window.location.pathname
    const segments = path.split('/').filter(Boolean)
    const id = segments[segments.length - 1]
    if (!id || id === 'view') {
        return { data: null as Article | null, error: null }
    }
    try {
        const data = await articlesApi.get(id)
        return { data, error: null }
    } catch (error) {
        return { data: null as Article | null, error: error as Error }
    }
}

function App() {
    const { data: article, loading, error } = useApi<Article | null>(fetchArticle)
    const [selectedImageIndex, setSelectedImageIndex] = useState(0)

    const { addItem } = useCart()

    const handleSelectImage = (index: number) => {
        setSelectedImageIndex(index)
    }

    const handleAddToCart = () => {
        if (!article) return
        addItem({
            articleId: article.Id,
            name: article.Name,
            price: article.Price,
            thumbnailUrl: article.ThumbnailUrl
        })
    }

    if (loading) return <Layout><p>Loading...</p></Layout>
    if (error) return <Layout><p>Error: {error.message}</p></Layout>
    if (!article) return <Layout><p>Article not found</p></Layout>

    return (
        <Layout>
            <section className={styles.details}>
                <div className={styles.imagesContainer}>
                    <div className={styles.display}>
                        {article.ImagesUrls.map((url, index) => (
                            <img
                                key={index}
                                className={styles.displayImage}
                                src={url}
                                alt={article.Name}
                                data-image_index={index}
                                data-selected={index === selectedImageIndex}
                                style={{ opacity: index === selectedImageIndex ? 1 : 0 }}
                            />
                        ))}
                    </div>

                    <div className={styles.miniatures}>
                        {article.ImagesUrls.map((url, index) => (
                            <button
                                key={index}
                                type="button"
                                onClick={() => handleSelectImage(index)}
                            >
                                <img src={url} alt={`${article.Name} thumbnail`} />
                            </button>
                        ))}
                    </div>
                </div>

                <div>
                    <h1 className={styles.name}>{article.Name}</h1>
                    <ArticleTags tags={article.Tags} />
                </div>

                <div>
                    {article.IsDeleted ? (
                        <span>Not available</span>
                    ) : (
                        <PriceDisplay price={article.Price} />
                    )}
                </div>

                {article.AvailableForTrade && (
                    <div>
                        <span>Available for trade!</span>
                        <a href="/wishlist">See the wishlist</a>
                    </div>
                )}

                <p className={styles.description}>{article.Description}</p>

                {!article.IsDeleted && (
                    <Button onClick={handleAddToCart}>Add to Cart</Button>
                )}
            </section>
        </Layout>
    )
}

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <App />
    </StrictMode>
)

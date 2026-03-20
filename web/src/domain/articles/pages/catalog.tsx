import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { Layout } from '@/shared/components'
import { CatalogList } from '../components'
import { catalogApi, type ArticleListItem } from '@/api'
import { useApi } from '@/shared/hooks'
import styles from './catalog.module.css'

function App() {
    return (
        <Layout>
            <Catalog />
        </Layout>
    )
}

async function fetchArticles() {
    try {
        const data = await catalogApi.list()
        return { data, error: null }
    } catch (error) {
        return { data: null as unknown as ArticleListItem[], error: error as Error }
    }
}

function Catalog() {
    const { data: articles, loading, error } = useApi<ArticleListItem[]>(fetchArticles)

    if (loading) return <p>Loading...</p>
    if (error) return <p>Error: {error.message}</p>

    return (
        <div className={styles.container}>
            <h1>Catalog</h1>
            {!articles || articles.length === 0 ? (
                <p>No articles available.</p>
            ) : (
                <CatalogList articles={articles} />
            )}
        </div>
    )
}

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <App />
    </StrictMode>
)

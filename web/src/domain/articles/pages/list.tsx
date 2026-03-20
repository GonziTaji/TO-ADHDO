import { useState, useEffect } from 'react'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { Layout, Button, Dialog } from '@/shared/components'
import { ArticleListItemComponent } from '../components'
import { articlesApi, type ArticleListItem } from '@/api'
import { useApi } from '@/shared/hooks'
import styles from './list.module.css'

function App() {
    return (
        <Layout>
            <ArticlesList />
        </Layout>
    )
}

async function fetchArticles() {
    try {
        const data = await articlesApi.list()
        return { data, error: null }
    } catch (error) {
        return { data: null as unknown as ArticleListItem[], error: error as Error }
    }
}

function ArticlesList() {
    const { data: initialArticles, loading, error } = useApi<ArticleListItem[]>(fetchArticles)
    const [articles, setArticles] = useState<ArticleListItem[]>([])
    const [deleteTarget, setDeleteTarget] = useState<string | null>(null)

    useEffect(() => {
        if (initialArticles) setArticles(initialArticles)
    }, [initialArticles])

    const handleNew = () => {
        window.location.href = '/articles/new'
    }

    const handleView = (id: string) => {
        window.location.href = `/catalog/${id}`
    }

    const handleEdit = (id: string) => {
        window.location.href = `/articles/${id}/edit`
    }

    const handleDelete = (id: string) => {
        setDeleteTarget(id)
    }

    const confirmDelete = async () => {
        if (!deleteTarget) return
        try {
            await articlesApi.delete(deleteTarget)
            setArticles((prev) => prev.filter((a) => a.Id !== deleteTarget))
            setDeleteTarget(null)
        } catch (err) {
            console.error('Failed to delete:', err)
        }
    }

    if (loading) return <p>Loading...</p>
    if (error) return <p>Error: {error.message}</p>

    return (
        <div id="articles_list" className={styles.container}>
            <h1>Articles list</h1>

            <Button onClick={handleNew}>New</Button>

            <ul className={styles.list}>
                {articles.map((article) => (
                    <ArticleListItemComponent
                        key={article.Id}
                        article={article}
                        onView={handleView}
                        onEdit={handleEdit}
                        onDelete={handleDelete}
                    />
                ))}
            </ul>

            <Dialog
                isOpen={deleteTarget !== null}
                onClose={() => setDeleteTarget(null)}
                title="Delete Article"
            >
                <p>Do you want to <strong>DELETE</strong> this article?</p>
                <div className={styles.dialogButtons}>
                    <Button variant="secondary" onClick={() => setDeleteTarget(null)}>
                        No, take me back
                    </Button>
                    <Button variant="danger" onClick={confirmDelete}>
                        Yes, delete the article
                    </Button>
                </div>
            </Dialog>
        </div>
    )
}

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <App />
    </StrictMode>
)

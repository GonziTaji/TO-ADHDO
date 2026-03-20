import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { Layout, Button } from '@/shared/components'
import { articlesApi, type Article } from '@/api'
import { useApi } from '@/shared/hooks'
import { useState, useEffect } from 'react'
import styles from './form.module.css'

function App() {
    return (
        <Layout>
            <ArticleForm />
        </Layout>
    )
}

async function fetchArticleById() {
    // Route: /articles/:id/edit -> id is the second-to-last segment
    // Route: /articles/new -> no id
    const segments = window.location.pathname.split('/').filter(Boolean)
    // segments for edit: ['articles', '<id>', 'edit']
    // segments for new:  ['articles', 'new']
    const id = segments.length >= 3 && segments[segments.length - 1] === 'edit'
        ? segments[segments.length - 2]
        : null

    if (!id) {
        return { data: null as Article | null, error: null }
    }
    try {
        const data = await articlesApi.get(id)
        return { data, error: null }
    } catch (error) {
        return { data: null as Article | null, error: error as Error }
    }
}

function ArticleForm() {
    const { data: fetchedArticle, loading: fetchLoading, error: fetchError } = useApi<Article | null>(fetchArticleById)

    const [article, setArticle] = useState<Partial<Article>>({
        Name: '',
        Description: '',
        Price: 0,
        Tags: []
    })
    const [submitLoading, setSubmitLoading] = useState(false)
    const [submitError, setSubmitError] = useState<string | null>(null)

    useEffect(() => {
        if (fetchedArticle) setArticle(fetchedArticle)
    }, [fetchedArticle])

    const handleChange = (field: keyof Article, value: unknown) => {
        setArticle((prev) => ({ ...prev, [field]: value }))
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        setSubmitError(null)
        setSubmitLoading(true)

        try {
            const segments = window.location.pathname.split('/').filter(Boolean)
            const id = segments.length >= 3 && segments[segments.length - 1] === 'edit'
                ? segments[segments.length - 2]
                : null

            if (id) {
                await articlesApi.update(id, article)
            } else {
                await articlesApi.create(article)
            }

            window.location.href = '/articles'
        } catch (err) {
            setSubmitError('Failed to save article')
            console.error(err)
        } finally {
            setSubmitLoading(false)
        }
    }

    if (fetchLoading) return <p>Loading...</p>
    if (fetchError) return <p>Error: {fetchError.message}</p>

    return (
        <form className={styles.form} onSubmit={handleSubmit}>
            <h1>{article.Name ? 'Edit Article' : 'New Article'}</h1>

            {submitError && <p className={styles.error}>{submitError}</p>}

            <label className={styles.field}>
                <span>Name</span>
                <input
                    type="text"
                    value={article.Name || ''}
                    onChange={(e) => handleChange('Name', e.target.value)}
                    required
                />
            </label>

            <label className={styles.field}>
                <span>Description</span>
                <textarea
                    value={article.Description || ''}
                    onChange={(e) => handleChange('Description', e.target.value)}
                    rows={4}
                />
            </label>

            <label className={styles.field}>
                <span>Price</span>
                <input
                    type="number"
                    value={article.Price || 0}
                    onChange={(e) => handleChange('Price', Number(e.target.value))}
                    min={0}
                />
            </label>

            <div className={styles.actions}>
                <Button type="submit" disabled={submitLoading}>
                    {submitLoading ? 'Saving...' : 'Save'}
                </Button>
                <Button type="button" variant="secondary" onClick={() => window.history.back()}>
                    Cancel
                </Button>
            </div>
        </form>
    )
}

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <App />
    </StrictMode>
)

import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { Layout } from '@/shared/components'
import { tagsApi, type Tag } from '@/api'
import { useApi } from '@/shared/hooks'
import styles from './Tags.module.css'

async function fetchTags() {
    try {
        const data = await tagsApi.list()
        return { data, error: null }
    } catch (error) {
        return { data: null as unknown as Tag[], error: error as Error }
    }
}

function Tags() {
    const { data: tags, loading, error } = useApi<Tag[]>(fetchTags)

    if (loading) return <Layout><p>Loading...</p></Layout>
    if (error) return <Layout><p>Error: {error.message}</p></Layout>

    return (
        <Layout>
            <div className={styles.container}>
                <h1>Tags</h1>
                <ul className={styles.list}>
                    {(tags ?? []).map((tag) => (
                        <li key={tag.Id} className={styles.item}>
                            <span>{tag.Name}</span>
                        </li>
                    ))}
                </ul>
            </div>
        </Layout>
    )
}

export { Tags as TagsPage }

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <Tags />
    </StrictMode>
)

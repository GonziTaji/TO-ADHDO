import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { Layout } from '@/shared/components'
import './styles/global.css'

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <Layout>
            <h1>Welcome to the Shop</h1>
            <p>Browse our catalog of articles.</p>
            <a href="/articles/catalog.html">View Catalog</a>
        </Layout>
    </StrictMode>
)

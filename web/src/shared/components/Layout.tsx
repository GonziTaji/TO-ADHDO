import { type ReactNode } from 'react'
import styles from './Layout.module.css'
import { CartProvider, useCart } from '@/shared/context/CartContext'

function Header() {
    const { items } = useCart()
    const itemCount = items.reduce((sum, item) => sum + item.quantity, 0)

    return (
        <header className={styles.header}>
            <nav className={styles.nav}>
                <a href="/">Home</a>
                <a href="/catalog">Catalog</a>
                <a href="/articles">Articles</a>
                <a href="/wishlist">Wishlist</a>
                <a href="/tags">Tags</a>
                <a href="/cart.html" className={styles.cartLink}>
                    Cart ({itemCount})
                </a>
            </nav>
        </header>
    )
}

function Footer() {
    return (
        <footer className={styles.footer}>
            <p>&copy; 2026 Personal Shop</p>
        </footer>
    )
}

interface LayoutProps {
    children: ReactNode
}

export function Layout({ children }: LayoutProps) {
    return (
        <CartProvider>
            <div className={styles.layout}>
                <Header />
                <main className={styles.main}>{children}</main>
                <Footer />
            </div>
        </CartProvider>
    )
}

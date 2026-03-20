import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { Layout, Button, PriceDisplay } from '@/shared/components'
import { useCart, type CartItem } from '@/shared/context'
import styles from './Cart.module.css'

function CartItemRow({ item, onUpdate, onRemove }: CartItemRowProps) {
    return (
        <li className={styles.item}>
            {item.thumbnailUrl && (
                <img src={item.thumbnailUrl} alt={item.name} className={styles.thumbnail} />
            )}
            <div className={styles.info}>
                <span className={styles.name}>{item.name}</span>
                <PriceDisplay price={item.price} />
            </div>
            <div className={styles.quantity}>
                <button onClick={() => onUpdate(item.id, item.quantity - 1)}>-</button>
                <span>{item.quantity}</span>
                <button onClick={() => onUpdate(item.id, item.quantity + 1)}>+</button>
            </div>
            <Button variant="danger" onClick={() => onRemove(item.id)}>Remove</Button>
        </li>
    )
}

interface CartItemRowProps {
    item: CartItem
    onUpdate: (id: string, quantity: number) => void
    onRemove: (id: string) => void
}

function App() {
    const { items, updateQuantity, removeItem, total, clearCart } = useCart()

    return (
        <Layout>
            <div className={styles.container}>
                <h1>Shopping Cart</h1>

                {items.length === 0 ? (
                    <p>Your cart is empty.</p>
                ) : (
                    <>
                        <ul className={styles.list}>
                            {items.map((item) => (
                                <CartItemRow
                                    key={item.id}
                                    item={item}
                                    onUpdate={updateQuantity}
                                    onRemove={removeItem}
                                />
                            ))}
                        </ul>

                        <div className={styles.summary}>
                            <div className={styles.total}>
                                <span>Total:</span>
                                <PriceDisplay price={total} />
                            </div>
                            <div className={styles.actions}>
                                <Button variant="secondary" onClick={clearCart}>Clear Cart</Button>
                                <Button>Checkout</Button>
                            </div>
                        </div>
                    </>
                )}
            </div>
        </Layout>
    )
}

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <App />
    </StrictMode>
)

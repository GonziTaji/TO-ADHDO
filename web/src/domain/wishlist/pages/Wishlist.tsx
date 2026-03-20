import { useState } from 'react'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { Layout, Tag } from '@/shared/components'
import type { WishlistResponse, WishlistItem } from '@/api'
import { wishlistApi } from '@/api'
import { useApi } from '@/shared/hooks'
import styles from './Wishlist.module.css'

async function fetchWishlist() {
    try {
        const data = await wishlistApi.list()
        return { data, error: null }
    } catch (error) {
        return { data: null as unknown as WishlistResponse, error: error as Error }
    }
}

interface WishlistItemRowProps {
    item: WishlistItem
}

function WishlistItemRow({ item }: WishlistItemRowProps) {
    return (
        <li>
            <span>{item.Name}</span>
        </li>
    )
}

interface WishlistFiltersProps {
    data: WishlistResponse
    selectedTags: string[]
    onTagToggle: (id: string) => void
    onSearch: (term: string) => void
    onPriceChange: (start: number, end: number) => void
    onApply: () => void
}

function WishlistFilters({ data, selectedTags, onTagToggle, onSearch, onPriceChange, onApply }: WishlistFiltersProps) {
    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        onApply()
    }

    return (
        <form className={styles.filtersForm} onSubmit={handleSubmit}>
            <fieldset>
                <legend>Search</legend>
                <div>
                    <input type="search" name="search" defaultValue={data.SearchTerm} onChange={(e) => onSearch(e.target.value)} />
                    <button type="submit">Search</button>
                </div>
            </fieldset>

            <fieldset>
                <legend>Price range</legend>
                <div>
                    <input
                        type="text"
                        pattern="\\d+"
                        inputMode="numeric"
                        name="price_start"
                        defaultValue={data.PriceSelectedRange.Start}
                        min={data.PriceRange.Start}
                        max={data.PriceRange.End}
                        onChange={(e) => onPriceChange(Number(e.target.value), data.PriceSelectedRange.End)}
                    />
                    -
                    <input
                        type="text"
                        pattern="\\d+"
                        inputMode="numeric"
                        name="price_end"
                        defaultValue={data.PriceSelectedRange.End}
                        min={data.PriceRange.Start}
                        max={data.PriceRange.End}
                        onChange={(e) => onPriceChange(data.PriceSelectedRange.Start, Number(e.target.value))}
                    />
                </div>
            </fieldset>

            <fieldset>
                <legend>Tags</legend>
                {data.TagsSelectOptions.map((tag) => (
                    <Tag
                        key={tag.Id}
                        name={tag.Name}
                        count={tag.Count}
                        selected={selectedTags.includes(tag.Id)}
                        disabled={tag.Count === 0}
                        onClick={() => onTagToggle(tag.Id)}
                    />
                ))}
            </fieldset>

            <button type="submit">Apply filters</button>
        </form>
    )
}

function Wishlist() {
    const { data, loading, error } = useApi<WishlistResponse>(fetchWishlist)
    const [selectedTags, setSelectedTags] = useState<string[]>([])

    if (loading) return <Layout><p>Loading...</p></Layout>
    if (error) return <Layout><p>Error: {error.message}</p></Layout>

    return (
        <Layout>
            <div className={styles.container}>
                <h1>My Wishlist</h1>

                {data && (
                    <WishlistFilters
                        data={data}
                        selectedTags={selectedTags}
                        onTagToggle={(id) => {
                            setSelectedTags((prev) =>
                                prev.includes(id) ? prev.filter((t) => t !== id) : [...prev, id]
                            )
                        }}
                        onSearch={(term) => console.log('Search:', term)}
                        onPriceChange={(start, end) => console.log('Price:', start, end)}
                        onApply={() => console.log('Apply filters')}
                    />
                )}

                <ul>
                    {(data?.Items ?? []).map((item) => (
                        <WishlistItemRow key={item.Id} item={item} />
                    ))}
                </ul>
            </div>
        </Layout>
    )
}

export { Wishlist as WishlistPage }

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <Wishlist />
    </StrictMode>
)

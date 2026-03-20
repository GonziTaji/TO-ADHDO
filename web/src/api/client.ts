const API_BASE = '/api'

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers
    },
    ...options
  })

  if (!response.ok) {
    throw new Error(`API Error: ${response.status} ${response.statusText}`)
  }

  return response.json()
}

export interface Tag {
  Id: string
  Name: string
}

export interface Price {
  Price: number
  Description: string
}

export interface ArticleCondition {
  Slug: string
  Label: string
}

export interface Article {
  Id: string
  Name: string
  Description: string
  Price: number
  ReferencePrice: number
  ThumbnailUrl: string
  ImagesUrls: string[]
  Tags: Tag[]
  Prices: Price[]
  IsDeleted: boolean
  AvailableForTrade: boolean
  Condition?: ArticleCondition
}

export interface ArticleListItem {
  Id: string
  Name: string
  Price: number
  ThumbnailUrl: string
  Tags: Tag[]
  Prices: Price[]
  AvailableForTrade: boolean
  Condition?: ArticleCondition
}

export const articlesApi = {
  /** Admin article list */
  list: () => request<ArticleListItem[]>('/articles'),
  /** Article detail (used by view page) */
  get: (id: string) => request<Article>(`/articles/${id}`),
  create: (data: Partial<Article>) =>
    request<{ Id: string }>('/articles', { method: 'POST', body: JSON.stringify(data) }),
  update: (id: string, data: Partial<Article>) =>
    request<{ Id: string }>(`/articles/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  delete: (id: string) => request<void>(`/articles/${id}`, { method: 'DELETE' }),
}

export const catalogApi = {
  /** Public catalog list with optional filters */
  list: (params?: { search?: string; tags?: string[]; trade?: boolean }) => {
    const qs = new URLSearchParams()
    if (params?.search) qs.set('search', params.search)
    if (params?.tags) params.tags.forEach((t) => qs.append('tags', t))
    if (params?.trade) qs.set('trade', 'true')
    const query = qs.toString()
    return request<ArticleListItem[]>(`/catalog${query ? '?' + query : ''}`)
  },
}

export interface WishlistItem {
  Id: string
  Name: string
  ObservedPrice: number
  Tags: Array<{ Id: string; Name: string }>
}

export interface WishlistFilter {
  SearchTerm: string
  PriceRange: { Start: number; End: number }
  PriceSelectedRange: { Start: number; End: number }
  TagsSelectOptions: Array<Tag & { Selected: boolean; Count: number }>
}

export interface WishlistResponse {
  Items: WishlistItem[]
  SearchTerm: string
  PriceRange: { Start: number; End: number }
  PriceSelectedRange: { Start: number; End: number }
  TagsSelectOptions: Array<Tag & { Selected: boolean; Count: number }>
}

export const wishlistApi = {
  /** Returns items + filter metadata in a single call */
  list: (params?: { search?: string; tags?: string[]; price_start?: number; price_end?: number }) => {
    const qs = new URLSearchParams()
    if (params?.search) qs.set('search', params.search)
    if (params?.tags) params.tags.forEach((t) => qs.append('tags', t))
    if (params?.price_start) qs.set('price_start', String(params.price_start))
    if (params?.price_end) qs.set('price_end', String(params.price_end))
    const query = qs.toString()
    return request<WishlistResponse>(`/wishlist${query ? '?' + query : ''}`)
  },
}

export const tagsApi = {
  list: () => request<Tag[]>('/tags'),
}

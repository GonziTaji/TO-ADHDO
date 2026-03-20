import { useState, useEffect } from 'react'

export interface ApiResult<T> {
  data: T
  error: Error | null
}

export interface UseApiResult<T> {
  data: T | null
  loading: boolean
  error: Error | null
}

export function useApi<T>(fn: () => Promise<ApiResult<T>>): UseApiResult<T> {
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    let cancelled = false

    fn().then((result) => {
      if (cancelled) return
      if (result.error) {
        setError(result.error)
      } else {
        setData(result.data)
      }
    }).finally(() => {
      if (!cancelled) setLoading(false)
    })

    return () => { cancelled = true }
  }, [])

  return { data, loading, error }
}

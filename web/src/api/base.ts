interface ApiClientOptions {
  headers?: Record<string, string>
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  body?: unknown

  noJsonParse?: boolean
}

export const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:4000'

export async function apiClient<T = unknown>(endpoint: string, options?: ApiClientOptions) {
  const noJsonParse = options?.noJsonParse ?? false

  const res = await fetch(baseUrl + endpoint, {
    headers: {
      'Content-Type': 'application/json',
      ...(noJsonParse ? {} : { Accept: 'application/json' }),
      ...(options?.headers || {}),
    },
    credentials: 'include',
    method: options?.method ?? 'GET',
    ...(options?.body ? { body: JSON.stringify(options.body) } : {}),
  })

  if (!res.ok && !noJsonParse) {
    const error = await res.json()
    console.error(endpoint, error, res)

    return Promise.reject(new Error(error?.message || error?.error || 'An error occurred'))
  }

  return !noJsonParse ? (res.json() as T) : null
}

export type PaginationType = {
  limit: number
  count: number
  page: number
  hasNextPage: boolean
  hasPreviousPage: boolean
}
export type PaginationResponse<T> = PaginationType & { data: T[] }

import { useQuery, useMutation, useQueryClient } from '@tanstack/solid-query'
import { type Accessor } from 'solid-js'
import { apiClient } from './base'
import type { CreateAircraftInput, Aircraft, PaginationRequestParams, PaginationResponse } from './types'

export const useAircraftsQuery = (params: Accessor<PaginationRequestParams>) => {
  return useQuery(() => ({
    queryKey: ['aircrafts', params()] as const,
    queryFn: () => {
      const { page, limit } = params()
      return apiClient<PaginationResponse<Aircraft>>(`/api/aircrafts?page=${page}&limit=${limit}`)
    },
  }))
}

export const useAircraftQuery = (id: Accessor<number | undefined>) => {
  return useQuery(() => ({
    queryKey: ['aircraft', id()] as const,
    queryFn: () => apiClient<Aircraft>(`/api/aircrafts/${id}`),
    enabled: id() !== undefined,
  }))
}

export const useCreateAircraftMutation = (options?: { onSuccess?: () => void; onError?: (error: Error) => void }) => {
  const queryClient = useQueryClient()

  return useMutation(() => ({
    mutationFn: (data: CreateAircraftInput) => {
      return apiClient<{ id: number }>('/api/aircrafts', { method: 'POST', body: data })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['aircrafts'] })
      options?.onSuccess?.()
    },
    onError: options?.onError,
  }))
}

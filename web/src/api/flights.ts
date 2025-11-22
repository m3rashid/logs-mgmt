import { useQuery, useMutation, useQueryClient } from '@tanstack/solid-query'
import { type Accessor } from 'solid-js'
import { apiClient, baseUrl } from './base'
import type { CreateFlightInput, Flight, FlightStatus, PaginationRequestParams, PaginationResponse } from './types'

export const useFlightsQuery = (params: Accessor<PaginationRequestParams>) => {
  return useQuery(() => ({
    queryKey: ['flights', params()] as const,
    queryFn: () => {
      const { page, limit } = params()
      return apiClient<PaginationResponse<Flight>>(`/api/flights?page=${page}&limit=${limit}&status=1`)
    },
  }))
}

export const useFlightQuery = (id: Accessor<number | undefined>) => {
  return useQuery(() => ({
    queryKey: ['flight', id()] as const,
    queryFn: () => apiClient<Flight>(`/api/flights/${id()}`),
    enabled: id() !== undefined,
  }))
}

export const useCreateFlightMutation = (options?: { onSuccess?: () => void; onError?: (error: Error) => void }) => {
  const queryClient = useQueryClient()

  return useMutation(() => ({
    mutationFn: (data: CreateFlightInput) => {
      return apiClient<{ id: number }>('/api/flights', { method: 'POST', body: data })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['flights'] })
      options?.onSuccess?.()
    },
    onError: options?.onError,
  }))
}

export type UpdateFlightInput = {
  status?: FlightStatus
  end_time?: Date
}
export const useUpdateFlightMutation = (options?: { onSuccess?: () => void; onError?: (error: Error) => void }) => {
  const queryClient = useQueryClient()

  return useMutation(() => ({
    mutationFn: ({ id, data }: { id: number; data: UpdateFlightInput }) => {
      return apiClient<Flight>(`/api/flights/${id}`, { method: 'PUT', body: data })
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['flights'] })
      queryClient.invalidateQueries({ queryKey: ['flight', variables.id] })
      options?.onSuccess?.()
    },
    onError: options?.onError,
  }))
}

export function getFlightRealtime(flightId: number) {
  return new EventSource(`${baseUrl}/api/events/${flightId}`)
}

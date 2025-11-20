import { apiClient, baseUrl, type PaginationResponse } from './base'
import type { Flight, CreateFlightInput } from './types'

export type GetFlightsParams = {
  page?: number
  limit?: number
}

export const getFlights = async (params: GetFlightsParams = {}) => {
  const { page = 1, limit = 10 } = params
  return apiClient<PaginationResponse<Flight>>(`/flights?page=${page}&limit=${limit}`)
}

export const getFlight = async (id: number) => {
  return apiClient<Flight>(`/flights/${id}`)
}

export const createFlight = async (data: CreateFlightInput) => {
  return apiClient<Flight>('/flights', {
    method: 'POST',
    body: data,
  })
}

export const getFlightRealtime = (flightId: number) => {
  return new EventSource(`${baseUrl}/flights/${flightId}/realtime`)
}

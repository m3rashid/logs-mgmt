import { apiClient, type PaginationResponse } from './base'
import type { Aircraft, CreateAircraftInput } from './types'

export type GetAircraftsParams = {
  page?: number
  limit?: number
}

export const getAircrafts = async (params: GetAircraftsParams = {}) => {
  const { page = 1, limit = 10 } = params
  return apiClient<PaginationResponse<Aircraft>>(`/aircrafts?page=${page}&limit=${limit}`)
}

export const getAircraft = async (id: number) => {
  return apiClient<Aircraft>(`/aircrafts/${id}`)
}

export const createAircraft = async (data: CreateAircraftInput) => {
  return apiClient<Aircraft>('/aircrafts', {
    method: 'POST',
    body: data,
  })
}

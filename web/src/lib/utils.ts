import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'
import type { AircraftStatus, FlightStatus } from '../api/types'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDate(date: string) {
  return new Date(date).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function getAircraftStatusLabel(status: AircraftStatus): string {
  const labels = ['Inactive', 'Active', 'Maintenance']
  return labels[status] || 'Unknown'
}

export function getFlightStatusLabel(status: FlightStatus): string {
  const labels = ['Scheduled', 'In Flight', 'Completed', 'Failed']
  return labels[status] || 'Unknown'
}

export function formatCoordinates(lat: number, lon: number): string {
  const latDir = lat >= 0 ? 'N' : 'S'
  const lonDir = lon >= 0 ? 'E' : 'W'
  return `${Math.abs(lat).toFixed(4)}°${latDir}, ${Math.abs(lon).toFixed(4)}°${lonDir}`
}

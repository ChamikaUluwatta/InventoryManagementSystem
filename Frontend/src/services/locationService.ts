import type { Location } from '@/types/location'
import { apiFetch } from '@/lib/api'

export const getAllLocations = async (): Promise<Location[]> => {
  return apiFetch<Location[]>('/locations')
}

export const getLocationById = async (id: string): Promise<Location> => {
  return apiFetch<Location>(`/locations/${id}`)
}

export const createLocation = async (
  location: Location,
): Promise<Location> => {
  return apiFetch<Location>('/locations', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(location),
  })
}

export const updateLocation = async (
  id: string,
  location: Partial<Location>,
): Promise<Location> => {
  return apiFetch<Location>(`/locations/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(location),
  })
}

export const deleteLocation = async (id: string): Promise<void> => {
  return apiFetch<void>(`/locations/${id}`, {
    method: 'DELETE',
  })
}

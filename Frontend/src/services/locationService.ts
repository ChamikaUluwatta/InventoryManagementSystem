import type { Location } from '@/types/location'

const API_BASE_URL = import.meta.env.VITE_API_URL;

if (!API_BASE_URL && import.meta.env.MODE === 'production') {
  throw new Error('VITE_API_URL environment variable is required for production build');
}

const API_BASE = API_BASE_URL || 'http://localhost:8080/api/v1'

export const getAllLocations = async (): Promise<Location[]> => {
  const response = await fetch(`${API_BASE}/locations`)
  if (!response.ok) {
    throw new Error('Failed to fetch locations')
  }
  return response.json()
}

export const getLocationById = async (id: string): Promise<Location> => {
  const response = await fetch(`${API_BASE}/locations/${id}`)
  if (!response.ok) {
    throw new Error('Failed to fetch location')
  }
  return response.json()
}

export const createLocation = async (
  location: Omit<Location, 'location_id'>,
): Promise<Location> => {
  const response = await fetch(`${API_BASE}/locations`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(location),
  })
  if (!response.ok) {
    throw new Error('Failed to create location')
  }
  return response.json()
}

export const updateLocation = async (
  id: string,
  location: Partial<Location>,
): Promise<Location> => {
  const response = await fetch(`${API_BASE}/locations/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(location),
  })
  if (!response.ok) {
    throw new Error('Failed to update location')
  }
  return response.json()
}

export const deleteLocation = async (id: string): Promise<void> => {
  const response = await fetch(`${API_BASE}/locations/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    throw new Error('Failed to delete location')
  }
}

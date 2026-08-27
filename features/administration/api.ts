import type { Media, Milestone, Project } from '../dashboard/types'

const request = async <T>(path: string, init?: RequestInit): Promise<T> => {
  const response = await fetch(path, { ...init, headers: { 'Content-Type': 'application/json', ...init?.headers } })
  if (!response.ok) {
    const body = await response.json().catch(() => ({ error: 'The local API request failed.' }))
    throw new Error(body.error || 'The local API request failed.')
  }
  return response.status === 204 ? (undefined as T) : response.json() as Promise<T>
}

export const listProjects = () => request<Project[]>('/api/projects', { cache: 'no-store' })
export const getProject = (id: string) => request<Project>(`/api/projects/${encodeURIComponent(id)}`, { cache: 'no-store' })
export const saveProject = (project: Partial<Project>, id?: string) => request<Project>(id ? `/api/projects/${encodeURIComponent(id)}` : '/api/projects', {
  method: id ? 'PUT' : 'POST',
  body: JSON.stringify(project),
})
export const addMedia = (projectID: string, media: Omit<Media, 'id' | 'project_id'>) => request<Media>(`/api/projects/${encodeURIComponent(projectID)}/media`, { method: 'POST', body: JSON.stringify(media) })
export const removeMedia = (projectID: string, mediaID: string) => request<void>(`/api/projects/${encodeURIComponent(projectID)}/media/${encodeURIComponent(mediaID)}`, { method: 'DELETE' })
export const addMilestone = (projectID: string, milestone: Omit<Milestone, 'id' | 'project_id'>) => request<Milestone>(`/api/projects/${encodeURIComponent(projectID)}/milestones`, { method: 'POST', body: JSON.stringify(milestone) })
export const updateMilestone = (projectID: string, milestoneID: string, milestone: Omit<Milestone, 'id' | 'project_id'>) => request<Project>(`/api/projects/${encodeURIComponent(projectID)}/milestones/${encodeURIComponent(milestoneID)}`, { method: 'PUT', body: JSON.stringify(milestone) })
export const removeMilestone = (projectID: string, milestoneID: string) => request<void>(`/api/projects/${encodeURIComponent(projectID)}/milestones/${encodeURIComponent(milestoneID)}`, { method: 'DELETE' })

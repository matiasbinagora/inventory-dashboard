import type { Project } from './types'

export function filterProjects(projects: Project[], search: string, technology: string, platform: string) {
  const query = search.trim().toLowerCase()
  return projects.filter((project) => {
    const matchesSearch = !query || [project.name, project.description ?? '', ...project.technologies].join(' ').toLowerCase().includes(query)
    const matchesTechnology = !technology || project.technologies.includes(technology)
    const matchesPlatform = !platform || project.agentic_platform === platform
    return matchesSearch && matchesTechnology && matchesPlatform
  })
}

export function countUpValue(target: number, elapsed: number, duration = 700) {
  if (target <= 0 || elapsed >= duration) return target
  return Math.floor(target * Math.min(1, Math.max(0, elapsed / duration)))
}

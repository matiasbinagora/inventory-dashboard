import { describe, expect, it } from 'vitest'
import { countUpValue, filterProjects } from './catalog'
import type { Project } from './types'

const projects: Project[] = [
  { id: 'one', name: 'Video Assistant', description: 'A media workflow', technologies: ['Go', 'React'], agentic_platform: 'OpenAI', links: [], media: [], milestones: [] },
  { id: 'two', name: 'Inference Lab', description: 'Model operations', technologies: ['Python'], agentic_platform: 'Claude', links: [], media: [], milestones: [] },
]

describe('catalog filters', () => {
  it('matches name, description, and technology without inventing fields', () => {
    expect(filterProjects(projects, 'media', '', '').map((project) => project.id)).toEqual(['one'])
    expect(filterProjects(projects, '', 'Python', '').map((project) => project.id)).toEqual(['two'])
  })
  it('combines technology and platform filters', () => {
    expect(filterProjects(projects, '', 'Go', 'Claude')).toEqual([])
    expect(filterProjects(projects, '', 'Go', 'OpenAI').map((project) => project.id)).toEqual(['one'])
  })
})

describe('deterministic count up', () => {
  it('always reaches the exact final value', () => {
    expect(countUpValue(17, 700)).toBe(17)
    expect(countUpValue(17, 900)).toBe(17)
    expect(countUpValue(17, -1)).toBe(0)
  })
})

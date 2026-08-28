import { describe, expect, it } from 'vitest'
import { fullSizeSource, imageMedia, localMediaURL, parseMetadataEntries, sortMilestones } from './detail'
import type { Project } from '../dashboard/types'

const project: Project = {
  id: 'atlas', name: 'Atlas', technologies: ['Go'], links: [],
  media: [
    { id: 'original', role: 'original', source: 'media/atlas/original.png', curated: true },
    { id: 'thumb', role: 'thumbnail', source: 'media/atlas/thumb.png', original_media_id: 'original', curated: true },
    { id: 'shot', role: 'screenshot', source: 'media/atlas/shot.png', curated: true },
  ],
  milestones: [
    { date: '2026-03-01', title: 'Second', description: 'Later', media_ids: [] },
    { date: '2026-01-01', title: 'First', description: 'Earlier', media_ids: [] },
  ],
}

describe('project detail data presentation', () => {
  it('sorts manual history chronologically without mutating the API response', () => {
    expect(sortMilestones(project.milestones).map((milestone) => milestone.title)).toEqual(['First', 'Second'])
    expect(project.milestones[0].title).toBe('Second')
  })
  it('resolves a thumbnail to its curated original and normalizes managed media paths', () => {
    expect(fullSizeSource(project.media[1], project.media)).toBe('media/atlas/original.png')
    expect(localMediaURL('media/atlas/original.png')).toBe('/media/atlas/original.png')
    expect(localMediaURL('https://cdn.example.test/demo.png')).toBe('https://cdn.example.test/demo.png')
  })
  it('keeps only image roles in the gallery', () => {
    expect(imageMedia(project).map((media) => media.role)).toEqual(['thumbnail', 'screenshot'])
  })
  it('parses named metadata links and preserves unexpected content', () => {
    expect(parseMetadataEntries(['Go: https://go.dev', 'Plain tool'])).toEqual([
      { primary: 'Go', urls: ['https://go.dev'] },
      { primary: 'Plain tool', urls: [] },
    ])
    expect(parseMetadataEntries('OpenAI — hosted docs: https://platform.openai.com/docs; Local runner')).toEqual([
      { primary: 'OpenAI — hosted docs', urls: ['https://platform.openai.com/docs'] },
      { primary: 'Local runner', urls: [] },
    ])
    expect(parseMetadataEntries('Python: https://python.org, Slack Bolt: https://slack.dev')).toEqual([
      { primary: 'Python', urls: ['https://python.org'] },
      { primary: 'Slack Bolt', urls: ['https://slack.dev'] },
    ])
  })
  it('deduplicates equivalent image sources without changing project media', () => {
    const duplicateProject = { ...project, media: [...project.media, { id: 'duplicate', role: 'screenshot' as const, source: 'media/atlas/original.png', curated: true }] }
    expect(imageMedia(duplicateProject).map((media) => media.id)).toEqual(['thumb', 'shot'])
    expect(duplicateProject.media).toHaveLength(4)
  })
})

import { describe, expect, it } from 'vitest'
import { fullSizeSource, imageMedia, localMediaURL, sortMilestones } from './detail'
import type { Project } from '../dashboard/types'

const project: Project = {
  id: 'atlas', name: 'Atlas', technologies: ['Go'], links: [],
  media: [
    { id: 'original', role: 'original', source: 'media/atlas/original.png' },
    { id: 'thumb', role: 'thumbnail', source: 'media/atlas/thumb.png', original_media_id: 'original' },
    { id: 'shot', role: 'screenshot', source: 'media/atlas/shot.png' },
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
    expect(imageMedia(project).map((media) => media.role)).toEqual(['original', 'thumbnail', 'screenshot'])
  })
})

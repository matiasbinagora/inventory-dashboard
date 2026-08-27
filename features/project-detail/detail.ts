import type { Media, Milestone, Project } from '../dashboard/types'

export function sortMilestones(milestones: Milestone[]) {
  return [...milestones].sort((a, b) => a.date.localeCompare(b.date))
}

export function localMediaURL(source: string) {
  return /^https:\/\//i.test(source) || source.startsWith('/') ? source : `/${source}`
}

export function fullSizeSource(media: Media, allMedia: Media[]) {
  if (media.original_media_id) {
    const original = allMedia.find((candidate) => candidate.id === media.original_media_id && candidate.role === 'original')
    if (original) return original.source
  }
  return media.role === 'original' ? media.source : media.source
}

export function imageMedia(project: Project) {
  return project.media.filter((media) => media.role === 'thumbnail' || media.role === 'screenshot' || media.role === 'original')
}

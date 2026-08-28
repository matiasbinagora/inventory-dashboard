import type { Media, Milestone, Project } from '../dashboard/types'
export { parseMetadataEntries } from '../shared/metadata'

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
  const images = project.media.filter((media) => media.role === 'thumbnail' || media.role === 'screenshot' || media.role === 'original')
  const originals = new Map(project.media.flatMap((media) => media.id && media.role === 'original' ? [[media.id, media.source] as const] : []))
  const representatives = new Map<string, Media>()
  for (const media of images) {
    const identity = media.original_media_id ? originals.get(media.original_media_id) || media.original_media_id : media.source
    const previous = representatives.get(identity)
    if (!previous || (media.role === 'thumbnail' && previous.role === 'original')) representatives.set(identity, media)
  }
  return images.filter((media) => representatives.get(media.original_media_id ? originals.get(media.original_media_id) || media.original_media_id : media.source) === media)
}

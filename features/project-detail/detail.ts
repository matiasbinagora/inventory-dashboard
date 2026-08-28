import type { Media, Milestone, Project } from '../dashboard/types'

export type MetadataEntry = { primary: string; urls: string[] }

const URL_PATTERN = /https:\/\/[^\s,;]+/gi

function cleanURL(value: string) {
  return value.replace(/[.)\]}]+$/g, '')
}

function parseMetadataSegment(value: string): MetadataEntry[] {
  const matches = [...value.matchAll(URL_PATTERN)]
  if (!matches.length) return [{ primary: value.trim(), urls: [] }]
  const entries = matches.map((match, index) => {
    const previousEnd = index === 0 ? 0 : (matches[index - 1].index || 0) + matches[index - 1][0].length
    const primary = value.slice(previousEnd, match.index || 0).replace(/^[\s:—–,;-]+|[\s:—–,;-]+$/g, '').trim()
    return { primary: primary || (index === 0 ? value.trim() : 'Reference'), urls: [cleanURL(match[0])] }
  })
  const last = matches[matches.length - 1]
  const trailing = value.slice((last.index || 0) + last[0].length).replace(/^[\s:—–,;-]+|[\s:—–,;-]+$/g, '').trim()
  return trailing ? [...entries, { primary: trailing, urls: [] }] : entries
}

export function parseMetadataEntries(value: string | string[]) {
  const values = Array.isArray(value) ? value : [value]
  return values.flatMap((item) => item.split(/\n|;\s*/).map((segment) => segment.trim()).filter(Boolean).flatMap(parseMetadataSegment))
}

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

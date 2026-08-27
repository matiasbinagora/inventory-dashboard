import type { Media } from '../dashboard/types'

export const thumbnailPayload = (media: Pick<Media, 'source' | 'alt_text' | 'caption'>): Omit<Media, 'id' | 'project_id'> => ({
  role: 'thumbnail', source: media.source, alt_text: media.alt_text || '', caption: media.caption || '', original_media_id: '',
})

import { describe, expect, it } from 'vitest'
import { thumbnailPayload } from './form'

describe('administration form helpers', () => {
  it('promotes a curated source without removing the original asset', () => {
    expect(thumbnailPayload({ source: 'media/original.png', alt_text: 'Original' })).toEqual({
      role: 'thumbnail', source: 'media/original.png', alt_text: 'Original', caption: '', original_media_id: '', curated: true,
    })
  })
})

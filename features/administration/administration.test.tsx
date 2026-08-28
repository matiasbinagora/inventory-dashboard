import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { thumbnailPayload } from './form'

vi.mock('./api', () => ({
  listProjects: vi.fn().mockResolvedValue([]),
  getProject: vi.fn(),
  saveProject: vi.fn(),
  addMedia: vi.fn(),
  removeMedia: vi.fn(),
  addMilestone: vi.fn(),
  updateMilestone: vi.fn(),
  removeMilestone: vi.fn(),
}))

import { Administration } from './administration'

afterEach(() => cleanup())

describe('administration form helpers', () => {
  it('promotes a curated source without removing the original asset', () => {
    expect(thumbnailPayload({ source: 'media/original.png', alt_text: 'Original' })).toEqual({
      role: 'thumbnail', source: 'media/original.png', alt_text: 'Original', caption: '', original_media_id: '', curated: true,
    })
  })

  it('renders an accessible, separated administration header', () => {
    render(<Administration />)

    const backLink = screen.getByRole('link', { name: '← Back to dashboard' })
    expect(backLink).toHaveAttribute('href', '/')
    expect(screen.getByRole('heading', { name: 'Inventory Administration', level: 1 })).toBeVisible()
    expect(backLink.parentElement).toHaveClass('admin-header-nav')
    expect(screen.getByText('LOCAL / ADMINISTRATION')).toBeVisible()
  })

  it('renders optional repository reference fields', () => {
    render(<Administration />)
    expect(screen.getByRole('textbox', { name: 'GitHub repository URL' })).toBeVisible()
    expect(screen.getByRole('textbox', { name: 'Trello backlog URL' })).toBeVisible()
  })
})

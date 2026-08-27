import { cleanup, render, screen, within } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { Dashboard } from './dashboard'

describe('dashboard sidebar navigation', () => {
  afterEach(() => {
    cleanup()
    vi.restoreAllMocks()
  })

  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: async () => [],
    }))
  })

  it('exposes Administration as an enabled keyboard-focusable link to /admin', () => {
    render(<Dashboard />)

    const navigation = screen.getByLabelText('Primary navigation')
    const administration = within(navigation).getByRole('link', { name: 'Administration' })

    expect(administration).toHaveAttribute('href', '/admin')
    expect(administration).not.toHaveAttribute('aria-disabled')
    expect(administration).not.toHaveAttribute('tabindex', '-1')
  })

  it('keeps Dashboard marked as the current sidebar destination', () => {
    render(<Dashboard />)

    const navigation = screen.getByLabelText('Primary navigation')
    expect(within(navigation).getByRole('link', { name: 'Dashboard' })).toHaveClass('active')
  })
})

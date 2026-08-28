import { cleanup, render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { ProjectDetail } from './project-detail'

const project = {
  id: 'atlas',
  name: 'Atlas',
  description: 'A readable project description.',
  technologies: ['Go', 'Next.js'],
  agentic_platform: 'OpenAI — https://platform.openai.com/docs',
  links: [
    { id: 'github', kind: 'github' as const, url: 'https://github.com/example/atlas', label: 'Source' },
    { id: 'trello', kind: 'trello' as const, url: 'https://trello.com/c/atlas', label: 'Plan' },
  ],
  media: [
    { id: 'image', role: 'screenshot' as const, source: 'media/atlas/screenshot.png', curated: true },
    { id: 'video', role: 'video' as const, source: 'https://video.example.test/atlas', caption: 'Demo' , curated: true },
  ],
  milestones: [],
}

afterEach(() => {
  cleanup()
  vi.restoreAllMocks()
})

describe('ProjectDetail presentation', () => {
  it('uses the editorial labels and semantic technology list', async () => {
    vi.spyOn(global, 'fetch').mockResolvedValue(new Response(JSON.stringify(project), { status: 200 }))

    render(<ProjectDetail projectID="atlas" />)
    await screen.findByRole('heading', { name: 'Atlas' })

    expect(screen.getByRole('heading', { name: 'Graphify Report' })).toBeVisible()
    expect(screen.getByRole('heading', { name: 'Demo Video' })).toBeVisible()
    expect(document.querySelector('.technology-list')).toContainElement(screen.getByText('Go'))
    expect(screen.getByText('Next.js').tagName).toBe('SPAN')
    expect(screen.getByRole('link', { name: 'https://platform.openai.com/docs' })).toHaveClass('metadata-link')
    expect(document.querySelector('.detail-hero p')).toHaveTextContent(project.description)
  })

  it('preserves public link destinations and makes platform references actionable', async () => {
    vi.spyOn(global, 'fetch').mockResolvedValue(new Response(JSON.stringify(project), { status: 200 }))

    render(<ProjectDetail projectID="atlas" />)
    await waitFor(() => expect(screen.getByRole('link', { name: 'Source' })).toBeVisible())

    expect(screen.getByRole('link', { name: 'Source' })).toHaveAttribute('href', project.links[0].url)
    expect(screen.getByRole('link', { name: 'Plan' })).toHaveAttribute('href', project.links[1].url)
    expect(screen.getByRole('link', { name: project.agentic_platform.split(' — ')[1] })).toHaveAttribute('href', project.agentic_platform.split(' — ')[1])
    expect(screen.getByRole('link', { name: 'Demo' })).toHaveAttribute('href', project.media[1].source)
  })

  it('renders the representative report at full content width while preserving detail links', async () => {
    vi.spyOn(global, 'fetch').mockResolvedValue(new Response(JSON.stringify(project), { status: 200 }))

    render(<ProjectDetail projectID="atlas" />)
    await screen.findByRole('heading', { name: 'Atlas' })

    const report = document.querySelector('.gallery-item')
    expect(report).toHaveClass('gallery-item')
    expect(report?.querySelector('img')).not.toBeNull()
    expect(screen.getByRole('link', { name: 'https://platform.openai.com/docs' })).toHaveAttribute('href', 'https://platform.openai.com/docs')
    expect(screen.getByRole('link', { name: 'Source' })).toHaveAttribute('href', project.links[0].url)
  })
})

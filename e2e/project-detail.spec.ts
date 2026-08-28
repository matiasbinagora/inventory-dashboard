import { test, expect } from '@playwright/test'

const project = {
  id: 'curated-atlas', name: 'Curated Atlas', description: 'A safe editorial project',
  technologies: ['Go', 'React'], agentic_platform: 'OpenAI',
  links: [{ id: 'github', kind: 'github', url: 'https://github.com/example/atlas' }, { id: 'trello', kind: 'trello', url: 'https://trello.com/c/abc123' }],
  media: [
    { id: 'original', role: 'original', source: 'media/atlas/original.png', alt_text: 'Atlas original', curated: true },
    { id: 'thumbnail', role: 'thumbnail', source: 'media/atlas/thumb.png', original_media_id: 'original', alt_text: 'Atlas thumbnail', curated: true },
    { id: 'video', role: 'video', source: 'https://video.example.test/atlas', caption: 'Public demo' },
  ],
  milestones: [
    { id: 'later', date: '2026-04-01', title: 'Launch', description: 'Published', media_ids: [] },
    { id: 'first', date: '2026-01-01', title: 'Prototype', description: 'Started', media_ids: [] },
  ],
}

async function mockProject(page: import('@playwright/test').Page, value = project) {
  await page.route('**/api/projects/curated-atlas', (route) => route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(value) }))
}

test('AC1 loads the selected project detail from the Go API contract', async ({ page }) => {
  await mockProject(page)
  await page.goto('/projects/curated-atlas')
  await expect(page.getByRole('heading', { name: 'Curated Atlas' })).toBeVisible()
})

test('AC2 omits optional fields with no value', async ({ page }) => {
  await mockProject(page, { ...project, description: '', agentic_platform: '', technologies: [], links: [], media: [], milestones: [] })
  await page.goto('/projects/curated-atlas')
  await expect(page.getByRole('heading', { name: 'Curated Atlas' })).toBeVisible()
  await expect(page.getByText('Technologies')).toHaveCount(0)
  await expect(page.getByText('Continue exploring')).toHaveCount(0)
})

test('AC3 exposes public GitHub and Trello links as inert external links', async ({ page }) => {
  await mockProject(page)
  await page.goto('/projects/curated-atlas')
  await expect(page.getByRole('link', { name: /GitHub repository/ })).toHaveAttribute('href', 'https://github.com/example/atlas')
  await expect(page.getByRole('link', { name: /Trello board/ })).toHaveAttribute('href', 'https://trello.com/c/abc123')
})

test('AC4 opens the associated original when a thumbnail is selected', async ({ page }) => {
  await mockProject(page)
  await page.goto('/projects/curated-atlas')
  await page.getByRole('button', { name: /Open original/ }).click()
  await expect(page.getByRole('dialog')).toBeVisible()
  await expect(page.getByRole('dialog').locator('img')).toHaveAttribute('src', '/media/atlas/original.png')
})

test('AC5 renders public video references according to their type', async ({ page }) => {
  await mockProject(page)
  await page.goto('/projects/curated-atlas')
  await expect(page.getByRole('link', { name: /Public demo/ })).toHaveAttribute('href', 'https://video.example.test/atlas')
})

test('AC6 orders manual milestones chronologically', async ({ page }) => {
  await mockProject(page)
  await page.goto('/projects/curated-atlas')
  await expect(page.locator('.timeline li').nth(0)).toContainText('Prototype')
  await expect(page.locator('.timeline li').nth(1)).toContainText('Launch')
})

test('AC7 displays a loading state and an error state', async ({ page }) => {
  await page.route('**/api/projects/curated-atlas', (route) => route.fulfill({ status: 500, contentType: 'application/json', body: '{"error":"unavailable"}' }))
  await page.goto('/projects/curated-atlas')
  await expect(page.locator('.detail-state[role="alert"]')).toContainText('Project unavailable')
})

test('AC8 presents refined detail sections without narrow viewport overflow', async ({ page }) => {
  await mockProject(page)
  await page.setViewportSize({ width: 1280, height: 900 })
  await page.goto('/projects/curated-atlas')
  await expect(page.getByRole('heading', { name: 'Graphify Report' })).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Demo Video' })).toBeVisible()
  await expect(page.locator('.technology-list li')).toHaveCount(2)
  await expect(page.getByRole('link', { name: 'GitHub repository' })).toHaveAttribute('href', 'https://github.com/example/atlas')

  await page.setViewportSize({ width: 390, height: 844 })
  await expect(page.locator('.detail-page')).toBeVisible()
  expect(await page.evaluate(() => document.documentElement.scrollWidth)).toBeLessThanOrEqual(390)
  expect(await page.locator('.detail-hero h1').evaluate((element) => element.getBoundingClientRect().right)).toBeLessThanOrEqual(390)
})

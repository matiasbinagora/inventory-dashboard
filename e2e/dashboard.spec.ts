import { test, expect } from '@playwright/test'

test('desktop dashboard exposes catalog controls and empty state', async ({ page }) => {
  await page.goto('/')
  await expect(page.getByRole('heading', { name: 'Project inventory' })).toBeVisible()
  await expect(page.getByLabel('Primary navigation')).toBeVisible()
  await expect(page.getByRole('textbox', { name: 'Search projects' })).toBeVisible()
  await expect(page.getByRole('button', { name: /Table/ })).toBeVisible()
  await expect(page.getByText(/catalog is ready|No projects match|Could not reach/)).toBeVisible()
})

test('sidebar remains accessible when collapsed', async ({ page }) => {
  await page.goto('/')
  const toggle = page.getByRole('button', { name: 'Collapse sidebar' })
  await toggle.click()
  await expect(page.getByRole('button', { name: 'Expand sidebar' })).toHaveAttribute('aria-expanded', 'false')
})

test('catalog switches presentation and filters API records', async ({ page }) => {
  await page.route('**/api/projects', async (route) => route.fulfill({
    status: 200,
    contentType: 'application/json',
    body: JSON.stringify([{ id: 'curated-one', name: 'Curated Atlas', description: 'A safe editorial project', technologies: ['Go'], agentic_platform: 'OpenAI', links: [], media: [], milestones: [] }]),
  }))
  await page.goto('/')
  await expect(page.getByRole('link', { name: /Curated Atlas/ })).toBeVisible()
  await page.getByRole('button', { name: /Table/ }).click()
  await expect(page.getByRole('table')).toBeVisible()
  await page.getByRole('textbox', { name: 'Search projects' }).fill('not found')
  await expect(page.getByText('No projects match these filters')).toBeVisible()
})

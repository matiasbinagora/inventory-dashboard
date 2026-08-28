import { test, expect } from '@playwright/test'

test('desktop dashboard exposes catalog controls and curated seed', async ({ page }) => {
  await page.goto('/')
  await expect(page.getByRole('heading', { name: 'Project inventory' })).toBeVisible()
  await expect(page.getByLabel('Primary navigation')).toBeVisible()
  await expect(page.getByRole('textbox', { name: 'Search projects' })).toBeVisible()
  await expect(page.getByRole('button', { name: /Table/ })).toBeVisible()
  await expect(page.getByText(/Slack Video Assistant|No projects match|Could not reach/)).toBeVisible()
})

test('sidebar remains accessible when collapsed', async ({ page }) => {
  await page.goto('/')
  const toggle = page.getByRole('button', { name: 'Collapse sidebar' })
  await toggle.click()
  await expect(page.getByRole('button', { name: 'Expand sidebar' })).toHaveAttribute('aria-expanded', 'false')
})

test('sidebar Administration link navigates to the administration screen', async ({ page }) => {
  await page.goto('/')
  await page.getByLabel('Primary navigation').getByRole('link', { name: 'Administration' }).click()
  await expect(page).toHaveURL(/\/admin$/)
  await expect(page.getByRole('heading', { name: 'Inventory Administration' })).toBeVisible()
})

test('sidebar Administration link activates with the keyboard', async ({ page }) => {
  await page.goto('/')
  const administration = page.getByLabel('Primary navigation').getByRole('link', { name: 'Administration' })
  await administration.focus()
  await page.keyboard.press('Enter')
  await expect(page).toHaveURL(/\/admin$/)
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

test('serves a local favicon without a missing-icon 404', async ({ page }) => {
  const missingIconResponses: string[] = []

  page.on('response', (response) => {
    if (response.status() === 404 && /\/(favicon|icon)\.(ico|svg)$/.test(new URL(response.url()).pathname)) {
      missingIconResponses.push(response.url())
    }
  })

  await page.goto('/')

  const iconHref = await page.locator('link[rel="icon"]').getAttribute('href')
  expect(iconHref).not.toBeNull()
  expect(new URL(iconHref!, page.url()).pathname).toBe('/icon.svg')
  expect((await page.request.get(new URL(iconHref!, page.url()).toString())).status()).toBe(200)
  expect(missingIconResponses).toEqual([])
})

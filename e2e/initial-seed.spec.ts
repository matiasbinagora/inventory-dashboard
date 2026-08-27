import { test, expect } from '@playwright/test'

test('initial curated projects are available in catalog, detail, and administration', async ({ page }) => {
  await page.goto('/')
  await expect(page.getByRole('link', { name: /Slack Video Assistant/ })).toBeVisible()
  await expect(page.getByRole('link', { name: /AWS Elemental Inference Smart Crop Demo/ })).toBeVisible()

  await page.getByRole('link', { name: /Slack Video Assistant/ }).click()
  await expect(page.getByRole('heading', { name: 'Slack Video Assistant' })).toBeVisible()
  await expect(page.getByText('Claude Agent SDK')).toBeVisible()

  await page.goto('/admin')
  await expect(page.getByRole('heading', { name: 'Inventory Administration' })).toBeVisible()
  await expect(page.getByLabel('Select project')).toContainText('Slack Video Assistant')
  await expect(page.getByLabel('Select project')).toContainText('AWS Elemental Inference Smart Crop Demo')
})

test('administration header stays readable on a narrow viewport and keeps back navigation', async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 })
  await page.goto('/admin')

  await expect(page.getByRole('heading', { name: 'Inventory Administration', level: 1 })).toBeVisible()
  const headerNav = page.locator('.admin-header-nav')
  const backLink = page.getByRole('link', { name: '← Back to dashboard' })
  const administrationLabel = page.getByText('LOCAL / ADMINISTRATION')
  await expect(headerNav).toBeVisible()
  const [backBox, labelBox] = await Promise.all([backLink.boundingBox(), administrationLabel.boundingBox()])
  expect(backBox).not.toBeNull()
  expect(labelBox).not.toBeNull()
  expect(labelBox!.x - (backBox!.x + backBox!.width)).toBeGreaterThanOrEqual(16)
  expect(await page.evaluate(() => document.documentElement.scrollWidth)).toBeLessThanOrEqual(390)

  await backLink.click()
  await expect(page).toHaveURL(/\/$/)
  await expect(page.getByRole('heading', { name: 'Project inventory' })).toBeVisible()
})

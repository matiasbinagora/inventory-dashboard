import { test, expect } from '@playwright/test'

test('initial curated projects are available in catalog, detail, and administration', async ({ page }) => {
  await page.goto('/')
  await expect(page.getByRole('link', { name: /Slack Video Assistant/ })).toBeVisible()
  await expect(page.getByRole('link', { name: /AWS Elemental Inference Smart Crop Demo/ })).toBeVisible()

  await page.getByRole('link', { name: /Slack Video Assistant/ }).click()
  await expect(page.getByRole('heading', { name: 'Slack Video Assistant' })).toBeVisible()
  await expect(page.getByText('Claude Agent SDK')).toBeVisible()

  await page.goto('/admin')
  await expect(page.getByRole('heading', { name: 'Curate the inventory' })).toBeVisible()
  await expect(page.getByLabel('Select project')).toContainText('Slack Video Assistant')
  await expect(page.getByLabel('Select project')).toContainText('AWS Elemental Inference Smart Crop Demo')
})

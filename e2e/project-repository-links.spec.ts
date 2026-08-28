import { test, expect } from '@playwright/test'

const repository = 'https://github.com/acme/inventory'
const backlog = 'https://trello.com/b/inventory'

test('administrator can create, edit, reload, and navigate dedicated project links', async ({ page }) => {
  let saved = { id: 'links-project', name: 'Links project', technologies: [], links: [], media: [], milestones: [], github_repository_url: repository, trello_backlog_url: backlog }
  await page.route('**/api/projects', async (route) => {
    if (route.request().method() === 'GET') return route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify([saved]) })
    saved = { ...saved, ...route.request().postDataJSON() }
    return route.fulfill({ status: 201, contentType: 'application/json', body: JSON.stringify(saved) })
  })
  await page.route('**/api/projects/links-project', async (route) => {
    if (route.request().method() === 'GET') return route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(saved) })
    saved = { ...saved, ...route.request().postDataJSON() }
    return route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(saved) })
  })

  await page.goto('/admin')
  await page.getByRole('textbox', { name: 'Project name' }).fill('Links project')
  await page.getByRole('textbox', { name: 'GitHub repository URL' }).fill(repository)
  await page.getByRole('textbox', { name: 'Trello backlog URL' }).fill(backlog)
  await page.getByRole('button', { name: 'Create project' }).click()
  await expect(page.getByText('Project created.')).toBeVisible()
  await page.reload()
  await expect(page.getByRole('textbox', { name: 'GitHub repository URL' })).toHaveValue(repository)
  await expect(page.getByRole('textbox', { name: 'Trello backlog URL' })).toHaveValue(backlog)

  await page.goto('/projects/links-project')
  await expect(page.getByRole('link', { name: 'GitHub repository' })).toHaveAttribute('href', repository)
  await expect(page.getByRole('link', { name: 'Trello backlog' })).toHaveAttribute('href', backlog)
})

test('administrator receives controlled API feedback for invalid URLs', async ({ page }) => {
  await page.route('**/api/projects', (route) => route.fulfill({ status: 400, contentType: 'application/json', body: JSON.stringify({ error: 'GitHub repository URL must be an HTTPS URL with a path' }) }))
  await page.goto('/admin')
  await page.getByRole('textbox', { name: 'Project name' }).fill('Invalid links')
  await page.getByRole('textbox', { name: 'GitHub repository URL' }).fill('http://github.com/acme/inventory')
  await page.getByRole('button', { name: 'Create project' }).click()
  await expect(page.locator('.form-error')).toContainText('HTTPS URL')
})

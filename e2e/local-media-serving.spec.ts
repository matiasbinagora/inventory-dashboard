import { test, expect } from '@playwright/test'

const image = Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=', 'base64')
const video = Buffer.from([0x1a, 0x45, 0xdf, 0xa3, 0x93, 0x42, 0x86, 0x81])

test('renders temporary local image and video media through managed paths', async ({ page }) => {
  await page.route('**/api/projects/media-fixture', (route) => route.fulfill({
    status: 200,
    contentType: 'application/json',
    body: JSON.stringify({
      id: 'media-fixture',
      name: 'Media fixture',
      technologies: [],
      links: [],
      media: [
        { id: 'image', role: 'original', source: 'media/fixture/image.png', alt_text: 'Temporary image', curated: true },
        { id: 'video', role: 'video', source: 'media/fixture/demo.webm', caption: 'Temporary video', curated: true },
      ],
      milestones: [],
    }),
  }))
  await page.route('**/media/fixture/image.png', (route) => route.fulfill({ status: 200, contentType: 'image/png', body: image }))
  await page.route('**/media/fixture/demo.webm', (route) => route.fulfill({ status: 200, contentType: 'video/webm', body: video }))

  const imageResponse = page.waitForResponse('**/media/fixture/image.png')
  await page.goto('/projects/media-fixture')
  await expect(page.getByAltText('Temporary image')).toBeVisible()
  await expect(page.locator('video source')).toHaveAttribute('src', '/media/fixture/demo.webm')
  expect((await imageResponse).status()).toBe(200)
})

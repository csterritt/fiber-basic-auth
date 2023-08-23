import { test, expect } from '@playwright/test'

test('test', async ({ page }) => {
  await page.goto('http://localhost:3000/protected')
  await expect(
    page.getByRole('heading', { name: 'Sign In Page' })
  ).toBeVisible()

  const text = await page.getByRole('alert').innerText()
  expect(
    text === 'There was a problem: You must be signed in to visit that page.'
  ).toBeTruthy()
})

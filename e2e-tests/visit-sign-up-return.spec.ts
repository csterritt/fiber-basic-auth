import { test, expect } from '@playwright/test'

test('test', async ({ page }) => {
  await page.goto('http://localhost:3000/')

  await expect(
    page.getByRole('heading', { name: 'A basic sign-in protected website.' })
  ).toBeVisible()

  await page.getByRole('link', { name: 'Sign up' }).click()
  await expect(
    page.getByRole('heading', { name: 'Sign Up Page' })
  ).toBeVisible()

  await page.getByRole('link', { name: 'Back to the main page.' }).click()

  await expect(
    page.getByRole('heading', { name: 'A basic sign-in protected website.' })
  ).toBeVisible()
})

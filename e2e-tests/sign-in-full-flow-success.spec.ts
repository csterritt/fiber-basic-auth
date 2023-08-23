import { test, expect } from '@playwright/test'

test('test', async ({ page }) => {
  await page.goto('http://localhost:3000/')
  await expect(
    page.getByRole('heading', { name: 'A basic sign-in protected website.' })
  ).toBeVisible()

  await page.getByRole('link', { name: 'Sign in' }).click()
  await expect(
    page.getByRole('heading', { name: 'Sign In Page' })
  ).toBeVisible()

  await page.getByPlaceholder('email').fill('x@yy.com')
  await page.getByRole('button', { name: 'Submit' }).click()

  await expect(
    page.getByRole('heading', {
      name: 'Enter your magic code here for email address x@yy.com.',
    })
  ).toBeVisible()

  await page.getByPlaceholder('code').fill('1234567')
  await page.getByRole('button', { name: 'Submit' }).click()

  const text = await page.getByRole('alert').innerText()
  expect(text === 'You are signed in.').toBeTruthy()

  await page.getByRole('button', { name: 'Sign out' }).click()

  await expect(page.getByRole('link', { name: 'Sign in' })).toBeVisible()
})

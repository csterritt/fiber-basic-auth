import { test, expect } from '@playwright/test'

test('test', async ({ page }) => {
  await page.goto('http://localhost:3000/')
  await expect(
    page.getByRole('heading', { name: 'A basic sign-in protected website.' })
  ).toBeVisible()

  await page.getByRole('link', { name: 'The protected page' }).click()
  await expect(
    page.getByRole('heading', { name: 'Sign In Page' })
  ).toBeVisible()
  await expect(
    page.getByText(
      'There was a problem: You must be signed in to visit that page.'
    )
  ).toBeVisible()

  await page.getByPlaceholder('email').fill('x@yy.com')
  await page.getByRole('button', { name: 'Submit' }).click()

  await expect(
    page.getByRole('heading', {
      name: 'Enter your magic code here for email address x@yy.com.',
    })
  ).toBeVisible()

  await page.getByPlaceholder('code').fill('1234')
  await page.getByRole('button', { name: 'Submit' }).click()

  await expect(
    page.getByRole('heading', {
      name: 'Enter your magic code here for email address x@yy.com.',
    })
  ).toBeVisible()
  await expect(
    page.getByText(
      'There was a problem: That code is incorrect or has expired.'
    )
  ).toBeVisible()
})

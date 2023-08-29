import { test } from '@playwright/test'

import { findHeading } from './support/finders'

test('test', async ({ page }) => {
  await page.goto('http://localhost:3000/')

  await findHeading(page, 'A basic sign-in protected website.')

  await page.getByRole('link', { name: 'Sign in' }).click()
  await findHeading(page, 'Sign In Page')

  await page.getByRole('link', { name: 'Back to the main page.' }).click()

  await findHeading(page, 'A basic sign-in protected website.')
})

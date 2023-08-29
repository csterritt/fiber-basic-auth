import { test } from '@playwright/test'

import { findHeading, findTextInRole } from './support/finders'

test('test', async ({ page }) => {
  await page.goto('http://localhost:3000/')
  await findHeading(page, 'A basic sign-in protected website.')

  await page.getByRole('link', { name: 'Sign up' }).click()
  await findHeading(page, 'Sign Up Page')

  await page.getByPlaceholder('email').fill('x@yy.com')
  await page.getByRole('button', { name: 'Submit' }).click()

  await findHeading(
    page,
    'Enter your magic code here for email address x@yy.com.'
  )

  await page.getByPlaceholder('code').fill('7654321')
  await page.getByRole('button', { name: 'Submit' }).click()

  await findHeading(page, 'Sign Up Page')
  await findTextInRole(
    page,
    'alert',
    'There was a problem: The code has expired, please try again.'
  )
})

import { test } from '@playwright/test'

import { findHeading, findTextInRole } from './support/finders'

test('test', async ({ page }) => {
  await page.goto('http://localhost:3000/')
  await findHeading(page, 'A basic sign-in protected website.')

  await page.getByRole('link', { name: 'The protected page' }).click()
  await findHeading(page, 'Sign In Page')
  await findTextInRole(
    page,
    'alert',
    'There was a problem: You must be signed in to visit that page.'
  )

  await page.getByPlaceholder('email').fill('x@yy.com')
  await page.getByRole('button', { name: 'Submit' }).click()

  await findHeading(
    page,
    'Enter your magic code here for email address x@yy.com.'
  )

  for (let i = 0; i < 3; i += 1) {
    await page.getByPlaceholder('code').fill('1234')
    await page.getByRole('button', { name: 'Submit' }).click()

    await findHeading(
      page,
      'Enter your magic code here for email address x@yy.com.'
    )
    await findTextInRole(
      page,
      'alert',
      'There was a problem: That code is incorrect.'
    )
  }

  await page.getByPlaceholder('code').fill('1234')
  await page.getByRole('button', { name: 'Submit' }).click()

  await findHeading(page, 'A basic sign-in protected website.')
  await findTextInRole(
    page,
    'alert',
    'There was a problem: The wrong code was given too many times.'
  )
})

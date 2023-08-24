import { test, expect } from '@playwright/test'

import { getPrefix, getCodeFromFileWithPrefix } from './support/code-access'
import { findHeading, findTextInRole } from './support/finders'

test('test', async ({ page }) => {
  await page.goto('http://localhost:3000/')
  await findHeading(page, 'A basic sign-in protected website.')

  await page.getByRole('link', { name: 'Sign in' }).click()
  await findHeading(page, 'Sign In Page')

  const prefix = getPrefix()
  await page.getByPlaceholder('email').fill(prefix + '+' + 'x@yy.com')
  await page.getByRole('button', { name: 'Submit' }).click()

  await findHeading(
    page,
    'Enter your magic code here for email address x@yy.com.'
  )

  const codeVal = getCodeFromFileWithPrefix(prefix)
  await page.getByPlaceholder('code').fill(codeVal)
  await page.getByRole('button', { name: 'Submit' }).click()

  await findTextInRole(page, 'alert', 'You are signed in.')

  await page.getByRole('button', { name: 'Sign out' }).click()

  await expect(page.getByRole('link', { name: 'Sign in' })).toBeVisible()
})

import { test } from '@playwright/test'

import { findHeading, findTextInRole } from './support/finders'

test('test', async ({ page }) => {
  await page.goto('http://localhost:3000/protected')
  await findHeading(page, 'Sign In Page')

  await findTextInRole(
    page,
    'alert',
    'There was a problem: You must be signed in to visit that page.'
  )
})

import { test, expect } from '@playwright/test'

import { findHeading } from './support/finders'

test('test', async ({ page }) => {
  await page.goto('http://localhost:3000/')
  await findHeading(page, 'A basic sign-in protected website.')
  const resp = await page.request.post(
    'http://localhost:3000/auth/submit-code',
    {
      data: { foo: 'bar' },
    }
  )
  expect(resp.status() === 303).toBeTruthy()
})

import * as fs from 'fs'
import { test, expect } from '@playwright/test'

const letters =
  'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'.split('')

const getRandomInt = (max: number) => Math.floor(Math.random() * max)

const randomString = () => {
  let res = ''
  for (let i = 0; i < 10; i += 1) {
    res += letters[getRandomInt(letters.length)]
  }

  return res
}

test('test', async ({ page }) => {
  await page.goto('http://localhost:3000/')
  await expect(
    page.getByRole('heading', { name: 'A basic sign-in protected website.' })
  ).toBeVisible()

  await page.getByRole('link', { name: 'Sign in' }).click()
  await expect(
    page.getByRole('heading', { name: 'Sign In Page' })
  ).toBeVisible()

  const prefix = randomString()
  await page.getByPlaceholder('email').fill(prefix + '+' + 'x@yy.com')
  await page.getByRole('button', { name: 'Submit' }).click()

  await expect(
    page.getByRole('heading', {
      name: 'Enter your magic code here for email address x@yy.com.',
    })
  ).toBeVisible()

  const filename = `/tmp/key-${prefix}.txt`
  const raw = fs.readFileSync(filename, 'utf8')
  fs.unlinkSync(filename)
  const codeVal = raw.toString()
  await page.getByPlaceholder('code').fill(codeVal)
  await page.getByRole('button', { name: 'Submit' }).click()

  const text = await page.getByRole('alert').innerText()
  expect(text === 'You are signed in.').toBeTruthy()

  await page.getByRole('button', { name: 'Sign out' }).click()

  await expect(page.getByRole('link', { name: 'Sign in' })).toBeVisible()
})

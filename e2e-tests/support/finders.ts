import { expect, Page } from '@playwright/test'

export const findHeading = async (page: Page, name: string) => {
  return expect(page.getByRole('heading', { name })).toBeVisible()
}

export const findTextInRole = async (page: Page, role: any, value: string) => {
  const text = await page.getByRole(role).innerText()
  return expect(text === value).toBeTruthy()
}

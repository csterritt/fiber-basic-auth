import * as fs from 'fs'
import { randomString } from './random-string'

export const getPrefix = () => randomString()

export const getCodeFromFileWithPrefix = (prefix: string): string => {
  const filename = `/tmp/key-${prefix}.txt`
  const raw = fs.readFileSync(filename, 'utf8')
  fs.unlinkSync(filename)
  return raw.toString()
}

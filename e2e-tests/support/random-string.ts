const letters =
  'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'.split('')

const getRandomInt = (max: number) => Math.floor(Math.random() * max)

export const randomString = () => {
  let res = ''
  for (let i = 0; i < 10; i += 1) {
    res += letters[getRandomInt(letters.length)]
  }

  return res
}

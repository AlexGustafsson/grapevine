export function useLocationPathPattern<const K extends string>(
  pattern: string,
  ...keys: K[]
): { [k in K]: string } | undefined {
  const match = new URLPattern(pattern, location.href).exec(location.href)
  if (!match) {
    return undefined
  }

  const result = {} as { [k in K]: string }
  for (const key of keys) {
    if (!match.pathname.groups[key]) {
      throw new Error('invalid pattern')
    }

    result[key] = match.pathname.groups[key]
  }

  return result
}

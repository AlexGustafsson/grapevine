export function useIsStandalone() {
  // Allow overrides for development / testing
  const url = new URL(location.href)

  switch (url.searchParams.get('standalone')) {
    case 'true':
      return true
    case 'false':
      return false
    default:
      return navigator.standalone
  }
}

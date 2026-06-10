export const REFRESH_INTERVAL_MS = 3000

export function startAutoRefresh(callback) {
	const timer = setInterval(callback, REFRESH_INTERVAL_MS)
	return () => clearInterval(timer)
}

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers || {}),
    },
    ...options,
  })

  const data = await response.json().catch(() => ({}))
  if (!response.ok) {
    throw new Error(data.error || 'Request failed')
  }
  return data
}

export function fetchPlatforms() {
  return request('/api/platforms')
}

export function generatePreviews(payload) {
  return request('/api/previews', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function publishContent(payload) {
  return request('/api/publish', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

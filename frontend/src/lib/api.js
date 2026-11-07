const BASE = '/api'

async function handle(res) {
  if (res.status === 204) return null
  const text = await res.text()
  const data = text ? safeParse(text) : null
  if (res.status === 401 && typeof window !== 'undefined') {
    window.dispatchEvent(new CustomEvent('lj:unauthorized'))
  }
  if (!res.ok) {
    const msg = (data && (data.error || data.message)) || res.statusText || 'Terjadi kesalahan'
    throw new Error(msg)
  }
  return data
}

function safeParse(text) {
  try {
    return JSON.parse(text)
  } catch {
    return null
  }
}

function query(params = {}) {
  const qs = new URLSearchParams()
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== null && value !== '') qs.set(key, value)
  }
  const s = qs.toString()
  return s ? `?${s}` : ''
}

function upload(url, file) {
  const form = new FormData()
  form.append('file', file)
  return fetch(url, { method: 'POST', body: form }).then(handle)
}

export const api = {
  authStatus: () => fetch(`${BASE}/auth/status`).then(handle),
  login: (password) =>
    fetch(`${BASE}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ password })
    }).then(handle),
  logout: () => fetch(`${BASE}/auth/logout`, { method: 'POST' }).then(handle),

  list: (params) => fetch(`${BASE}/kotoba${query(params)}`).then(handle),
  listKotobaLessons: () => fetch(`${BASE}/kotoba/lessons`).then(handle),
  get: (id) => fetch(`${BASE}/kotoba/${id}`).then(handle),
  create: (payload) =>
    fetch(`${BASE}/kotoba`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    }).then(handle),
  update: (id, payload) =>
    fetch(`${BASE}/kotoba/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    }).then(handle),
  remove: (id) => fetch(`${BASE}/kotoba/${id}`, { method: 'DELETE' }).then(handle),
  bulkKotoba: (file) => upload(`${BASE}/kotoba/bulk`, file),
  bulkRemoveKotoba: (ids) =>
    fetch(`${BASE}/kotoba/bulk-delete`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ids })
    }).then(handle),

  listKanji: (params) => fetch(`${BASE}/kanji${query(params)}`).then(handle),
  getKanji: (id) => fetch(`${BASE}/kanji/${id}`).then(handle),
  createKanji: (payload) =>
    fetch(`${BASE}/kanji`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    }).then(handle),
  updateKanji: (id, payload) =>
    fetch(`${BASE}/kanji/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    }).then(handle),
  removeKanji: (id) => fetch(`${BASE}/kanji/${id}`, { method: 'DELETE' }).then(handle),
  bulkKanji: (file) => upload(`${BASE}/kanji/bulk`, file),
  bulkRemoveKanji: (ids) =>
    fetch(`${BASE}/kanji/bulk-delete`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ids })
    }).then(handle),

  listBunpou: (params) => fetch(`${BASE}/bunpou${query(params)}`).then(handle),
  getBunpou: (id) => fetch(`${BASE}/bunpou/${id}`).then(handle),
  createBunpou: (payload) =>
    fetch(`${BASE}/bunpou`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    }).then(handle),
  updateBunpou: (id, payload) =>
    fetch(`${BASE}/bunpou/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    }).then(handle),
  removeBunpou: (id) => fetch(`${BASE}/bunpou/${id}`, { method: 'DELETE' }).then(handle),
  bulkRemoveBunpou: (ids) =>
    fetch(`${BASE}/bunpou/bulk-delete`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ids })
    }).then(handle),

  listDokkaiLessons: (type) =>
    fetch(`${BASE}/dokkai/${encodeURIComponent(type)}/lessons`).then(handle),
  getDokkai: (type, lesson) =>
    fetch(`${BASE}/dokkai/${type}/${encodeURIComponent(lesson)}`).then(handle),
  bulkDokkai: (type, file) => {
    const form = new FormData()
    form.append('type', type)
    form.append('file', file)
    return fetch(`${BASE}/dokkai/bulk`, { method: 'POST', body: form }).then(handle)
  },

  listDokkaiLevels: () => fetch(`${BASE}/dokkai/level`).then(handle),
  listDokkaiLevel: (level) =>
    fetch(`${BASE}/dokkai/level/${encodeURIComponent(level)}`).then(handle),
  bulkDokkaiLevel: (file, level) => {
    const form = new FormData()
    form.append('file', file)
    if (level) form.append('level', level)
    return fetch(`${BASE}/dokkai/level/bulk`, { method: 'POST', body: form }).then(handle)
  },
  updateDokkaiLevel: (level, id, payload) =>
    fetch(`${BASE}/dokkai/level/${encodeURIComponent(level)}/${encodeURIComponent(id)}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    }).then(handle),
  removeDokkaiLevel: (level, id) =>
    fetch(`${BASE}/dokkai/level/${encodeURIComponent(level)}/${encodeURIComponent(id)}`, {
      method: 'DELETE'
    }).then(handle)
}

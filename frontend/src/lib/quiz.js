export function shuffle(arr) {
  const a = [...arr]
  for (let i = a.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[a[i], a[j]] = [a[j], a[i]]
  }
  return a
}

export function orderById(items) {
  return [...items].sort((a, b) => {
    const na = Number(String(a?.id ?? '').replace(/\D/g, ''))
    const nb = Number(String(b?.id ?? '').replace(/\D/g, ''))
    if (!Number.isNaN(na) && !Number.isNaN(nb) && na !== nb) return na - nb
    return String(a?.id ?? '').localeCompare(String(b?.id ?? ''))
  })
}

export function neighborItems(word, items, count = 3) {
  const ordered = orderById(items)
  const idx = ordered.findIndex((w) => w.id === word.id)
  if (idx === -1) {
    return shuffle(items.filter((w) => w.id !== word.id)).slice(0, count)
  }
  const picked = []
  for (let d = 1; picked.length < count && d < ordered.length; d++) {
    const before = ordered[idx - d]
    const after = ordered[idx + d]
    if (before && before.id !== word.id) picked.push(before)
    if (after && picked.length < count) picked.push(after)
  }
  if (picked.length < count) {
    const chosen = new Set([word.id, ...picked.map((w) => w.id)])
    const rest = shuffle(
      items.filter((w) => !chosen.has(w.id))
    )
    for (const w of rest) {
      if (picked.length >= count) break
      picked.push(w)
    }
  }
  return picked
}

export function uniqueSortedLessons(items) {
  return [...new Set(items.map((i) => String(i.lesson || '')).filter(Boolean))].sort((a, b) => {
    const na = Number(a)
    const nb = Number(b)
    if (!Number.isNaN(na) && !Number.isNaN(nb) && na !== nb) return na - nb
    return a.localeCompare(b)
  })
}

const NUMBER_PATTERN = /[0-9０-９〇零一二三四五六七八九十百千万億―—]/

export function isDateWord(word) {
  const fields = [
    word?.kanji,
    word?.furigana,
    ...(word?.terjemahan || []),
  ]
  return fields.some((field) => {
    const value = String(field ?? '')
    return NUMBER_PATTERN.test(value) || value.toLowerCase().includes('tanggal')
  })
}

export function normalizeAnswer(value) {
  return String(value ?? '')
    .replace(/[（(][^）)]*[）)]/g, '')
    .trim()
    .toLowerCase()
    .replace(/\s+/g, '')
    .replace(/[.。,，、!！?？'"`’“”\-ー…]/g, '')
}

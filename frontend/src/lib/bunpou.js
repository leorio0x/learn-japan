import { parseParagraph } from './paragraph.js'

function lessonNumber(value) {
  const match = String(value ?? '').match(/\d+/)
  return match ? Number(match[0]) : Number.POSITIVE_INFINITY
}

export function sortLessons(items) {
  const unique = [
    ...new Set(
      items.map((item) => String(item?.lesson ?? '').trim()).filter(Boolean),
    ),
  ]
  return unique.sort((a, b) => {
    const na = lessonNumber(a)
    const nb = lessonNumber(b)
    if (na !== nb) return na - nb
    return a.localeCompare(b)
  })
}

export function lessonLabel(value) {
  const text = String(value ?? '').trim()
  if (!text) return ''
  return /lesson/i.test(text) ? text : `Lesson ${text}`
}

export function parseRumus(value) {
  const source = String(value ?? '')
  const parts = []
  const pattern = /\[([^\]]+)\]/g
  let last = 0
  let match

  while ((match = pattern.exec(source))) {
    if (match.index > last) {
      parts.push({ type: 'text', value: source.slice(last, match.index) })
    }

    const raw = match[1]
    const separator = raw.indexOf(':')
    const label = (separator >= 0 ? raw.slice(0, separator) : raw).trim()
    const categoriesRaw = separator >= 0 ? raw.slice(separator + 1) : ''
    const categories = categoriesRaw
      .split('|')
      .map((item) => item.trim())
      .filter(Boolean)

    parts.push({
      type: 'slot',
      value: label || '～',
      categories,
    })
    last = pattern.lastIndex
  }

  if (last < source.length) {
    parts.push({ type: 'text', value: source.slice(last) })
  }

  return parts
}

export function plainText(value, mode = 'kanji') {
  return parseParagraph(value)
    .map((part) => {
      if (part.type === 'text') return part.value
      if (mode === 'kana') return part.kana || part.kanji
      return part.kanji || part.kana
    })
    .join('')
    .trim()
}

export function targetAnswers(kalimat) {
  const kanji = plainText(kalimat, 'kanji')
  const kana = plainText(kalimat, 'kana')
  return [...new Set([kanji, kana].filter(Boolean))]
}

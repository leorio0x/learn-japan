export function stripJapaneseSpaces(value) {
  return String(value ?? '').replace(/[^\S\r\n]+/g, '')
}

export function parseParagraph(text) {
  const source = String(text ?? '')
  const re = /<([^<>]+)>/g
  const parts = []
  let last = 0
  let match

  while ((match = re.exec(source))) {
    if (match.index > last) {
      const value = stripJapaneseSpaces(source.slice(last, match.index))
      if (value) parts.push({ type: 'text', value })
    }

    const raw = match[1]
    const first = raw.indexOf(':')
    const second = first >= 0 ? raw.indexOf(':', first + 1) : -1
    let kanji = raw
    let kana = ''
    let arti = ''
    if (second >= 0) {
      kanji = raw.slice(0, first)
      kana = raw.slice(first + 1, second)
      arti = raw.slice(second + 1)
    } else if (first >= 0) {
      kanji = raw.slice(0, first)
      kana = raw.slice(first + 1)
    }

    parts.push({ type: 'word', kanji, kana, arti })
    last = re.lastIndex
  }

  if (last < source.length) {
    const value = stripJapaneseSpaces(source.slice(last))
    if (value) parts.push({ type: 'text', value })
  }

  return parts
}

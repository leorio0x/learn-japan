const STORAGE_KEY = 'lj:kanji-marked'

function readIds() {
  if (typeof window === 'undefined') return []
  try {
    const raw = window.sessionStorage.getItem(STORAGE_KEY)
    if (!raw) return []
    const parsed = JSON.parse(raw)
    if (!Array.isArray(parsed)) return []
    return [...new Set(parsed.map((id) => String(id)).filter(Boolean))]
  } catch {
    return []
  }
}

let ids = $state(readIds())
let quizRequested = $state(false)

function persist(next) {
  if (typeof window === 'undefined') return
  try {
    window.sessionStorage.setItem(STORAGE_KEY, JSON.stringify(next))
  } catch {
    // penyimpanan tidak tersedia (mis. private mode), abaikan
  }
}

function apply(next) {
  ids = next
  persist(next)
}

export const kanjiMark = {
  get ids() {
    return ids
  },
  get count() {
    return ids.length
  },
  has(id) {
    return ids.includes(String(id))
  },
  toggle(id) {
    const key = String(id)
    if (!key) return
    apply(ids.includes(key) ? ids.filter((x) => x !== key) : [...ids, key])
  },
  remove(id) {
    const key = String(id)
    if (!ids.includes(key)) return
    apply(ids.filter((x) => x !== key))
  },
  clear() {
    if (!ids.length) return
    apply([])
  },
  requestQuiz() {
    quizRequested = true
  },
  consumeQuizRequest() {
    const requested = quizRequested
    quizRequested = false
    return requested
  },
}

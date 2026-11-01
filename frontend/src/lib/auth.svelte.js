import { api } from './api.js'

let authenticated = $state(false)
let checking = $state(true)

export const auth = {
  get authenticated() {
    return authenticated
  },
  get checking() {
    return checking
  },
  async refresh() {
    checking = true
    try {
      const res = await api.authStatus()
      authenticated = !!res?.authenticated
    } catch {
      authenticated = false
    } finally {
      checking = false
    }
  },
  async login(password) {
    const res = await api.login(password)
    authenticated = !!res?.authenticated
    return res
  },
  async logout() {
    try {
      await api.logout()
    } catch {
      // abaikan, tetap kunci di sisi klien
    }
    authenticated = false
  },
  clear() {
    authenticated = false
  }
}

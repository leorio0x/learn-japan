<script>
  import { onMount } from 'svelte'
  import Navbar from './components/Navbar.svelte'
  import Kotoba from './pages/Kotoba.svelte'
  import Kanji from './pages/Kanji.svelte'
  import Kaiwa from './pages/Kaiwa.svelte'
  import Quiz from './pages/Quiz.svelte'
  import QuizKanji from './pages/QuizKanji.svelte'
  import QuizKaiwa from './pages/QuizKaiwa.svelte'
  import Crud from './pages/Crud.svelte'
  import CrudKanji from './pages/CrudKanji.svelte'
  import CrudKaiwa from './pages/CrudKaiwa.svelte'
  import CrudDokkai from './pages/CrudDokkai.svelte'
  import CrudGate from './components/CrudGate.svelte'
  import DokkaiView from './components/DokkaiView.svelte'
  import Dokkai from './pages/Dokkai.svelte'
  import Admin from './pages/Admin.svelte'
  import { auth } from './lib/auth.svelte.js'

  const routes = {
    '/': 'kotoba',
    '/kotoba': 'kotoba',
    '/kanji': 'kanji',
    '/kaiwa': 'kaiwa',
    '/quiz/kotoba': 'quiz-kotoba',
    '/quiz/kanji': 'quiz-kanji',
    '/quiz/kaiwa': 'quiz-kaiwa',
    '/dokkai': 'dokkai',
    '/dokkai/kotoba': 'dokkai-kotoba',
    '/dokkai/kanji': 'dokkai-kanji',
    '/kelola/kotoba': 'crud-kotoba',
    '/kelola/kanji': 'crud-kanji',
    '/kelola/kaiwa': 'crud-kaiwa',
    '/kelola/dokkai': 'crud-dokkai',
    '/ansatsunoadmin': 'admin'
  }

  const paths = Object.fromEntries(
    Object.entries(routes).map(([path, id]) => [id, path])
  )

  function resolve(pathname) {
    const path = (pathname || '/').replace(/\/+$/, '') || '/'
    return routes[path] || 'kotoba'
  }

  let page = $state(
    typeof window === 'undefined' ? 'kotoba' : resolve(window.location.pathname)
  )

  function navigate(next) {
    page = next
    if (typeof window === 'undefined') return
    const path = paths[next] || '/'
    if (window.location.pathname !== path) {
      window.history.pushState({ page: next }, '', path)
    }
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  onMount(() => {
    auth.refresh()

    const onPopState = () => {
      page = resolve(window.location.pathname)
    }
    const onUnauthorized = () => auth.clear()

    window.addEventListener('popstate', onPopState)
    window.addEventListener('lj:unauthorized', onUnauthorized)
    return () => {
      window.removeEventListener('popstate', onPopState)
      window.removeEventListener('lj:unauthorized', onUnauthorized)
    }
  })
</script>

<div class="flex min-h-screen flex-col bg-cream">
  <Navbar {page} onNavigate={navigate} />

  <main class="mx-auto w-full max-w-5xl flex-1 px-4 py-6 sm:px-6 sm:py-8">
    {#if page === 'kotoba'}
      <Kotoba onNavigate={navigate} />
    {:else if page === 'kanji'}
      <Kanji onNavigate={navigate} />
    {:else if page === 'kaiwa'}
      <Kaiwa onNavigate={navigate} />
    {:else if page === 'quiz-kotoba'}
      <Quiz />
    {:else if page === 'quiz-kanji'}
      <QuizKanji />
    {:else if page === 'quiz-kaiwa'}
      <QuizKaiwa />
    {:else if page === 'dokkai-kotoba'}
      <DokkaiView type="kotoba" onNavigate={navigate} />
    {:else if page === 'dokkai-kanji'}
      <DokkaiView type="kanji" onNavigate={navigate} />
    {:else if page === 'dokkai'}
      <Dokkai onNavigate={navigate} />
    {:else if page === 'crud-kotoba'}
      <CrudGate><Crud /></CrudGate>
    {:else if page === 'crud-kanji'}
      <CrudGate><CrudKanji /></CrudGate>
    {:else if page === 'crud-kaiwa'}
      <CrudGate><CrudKaiwa /></CrudGate>
    {:else if page === 'crud-dokkai'}
      <CrudGate><CrudDokkai onNavigate={navigate} /></CrudGate>
    {:else if page === 'admin'}
      <Admin onNavigate={navigate} />
    {:else}
      <Kotoba onNavigate={navigate} />
    {/if}
  </main>

  <footer class="border-t border-navy/10 py-5 text-center text-xs text-navy/50">
    Learn Japan &middot; Quiz Kotoba, Kanji, Kaiwa &amp; Dokkai
  </footer>
</div>

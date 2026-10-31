<script>
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
  import CrudGate from './components/CrudGate.svelte'

  let page = $state('kotoba')
</script>

<div class="flex min-h-screen flex-col bg-cream">
  <Navbar {page} onNavigate={(p) => (page = p)} />

  <main class="mx-auto w-full max-w-5xl flex-1 px-4 py-6 sm:px-6 sm:py-8">
    {#if page === 'kotoba'}
      <Kotoba />
    {:else if page === 'kanji'}
      <Kanji />
    {:else if page === 'kaiwa'}
      <Kaiwa />
    {:else if page === 'quiz-kotoba'}
      <Quiz />
    {:else if page === 'quiz-kanji'}
      <QuizKanji />
    {:else if page === 'quiz-kaiwa'}
      <QuizKaiwa />
    {:else if page === 'crud-kotoba'}
      <CrudGate><Crud /></CrudGate>
    {:else if page === 'crud-kanji'}
      <CrudGate><CrudKanji /></CrudGate>
    {:else}
      <CrudGate><CrudKaiwa /></CrudGate>
    {/if}
  </main>

  <footer class="border-t border-navy/10 py-5 text-center text-xs text-navy/50">
    Learn Japan &middot; Quiz Kotoba, Kanji &amp; Kaiwa
  </footer>
</div>

<script>
  let { kanji = '', kana = '', className = '', hideKanji = false } = $props()

  const hasKanji = $derived(!!(kanji && kanji.trim() !== ''))
  const hasKana = $derived(!!(kana && kana.trim() !== ''))
  const showRuby = $derived(hasKanji && hasKana && kana.trim() !== kanji.trim())
</script>

{#if hideKanji}
  <span class={className}>{hasKana ? kana : kanji}</span>
{:else if showRuby}
  <ruby class="ruby-jp {className}">{kanji}<rt>{kana}</rt></ruby>
{:else}
  <span class={className}>{kanji || kana}</span>
{/if}

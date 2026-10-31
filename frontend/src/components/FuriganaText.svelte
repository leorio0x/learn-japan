<script>
  let { text = '' } = $props()

  const parts = $derived(parse(text))

  function parse(value) {
    const source = String(value ?? '')
    const pattern = /<([^<>:]+):([^<>]+)>/g
    const result = []
    let last = 0
    let match

    while ((match = pattern.exec(source))) {
      if (match.index > last) {
        result.push({ type: 'text', value: source.slice(last, match.index) })
      }
      result.push({ type: 'ruby', kanji: match[1], kana: match[2] })
      last = pattern.lastIndex
    }

    if (last < source.length) {
      result.push({ type: 'text', value: source.slice(last) })
    }

    return result
  }
</script>

{#each parts as part}
  {#if part.type === 'ruby'}
    <ruby class="ruby-jp">{part.kanji}<rt>{part.kana}</rt></ruby>
  {:else}
    {part.value}
  {/if}
{/each}

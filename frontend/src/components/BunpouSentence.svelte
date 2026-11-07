<script>
    import { parseParagraph } from "../lib/paragraph.js";
    import DokkaiWord from "./DokkaiWord.svelte";

    let {
        kalimat = "",
        terjemahan = "",
        hideKanji = false,
        showTerjemahan = true,
        className = "",
    } = $props();

    const parts = $derived(parseParagraph(kalimat));
</script>

<span class="block">
    <span class="font-jp {className}">{#each parts as part}{#if part.type === "text"}{part.value}{:else}<DokkaiWord kanji={part.kanji} kana={part.kana} arti={part.arti} {hideKanji} />{/if}{/each}</span>
    {#if showTerjemahan && terjemahan}
        <span class="mt-0.5 block text-sm text-navy/55">{terjemahan}</span>
    {/if}
</span>

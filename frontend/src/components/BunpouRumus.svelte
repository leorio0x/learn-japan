<script>
    import { parseRumus } from "../lib/bunpou.js";
    import RumusSlot from "./RumusSlot.svelte";

    let { rumus = "", className = "" } = $props();

    const parts = $derived(parseRumus(rumus));
</script>

<span class="font-jp {className}">
    {#each parts as part}
        {#if part.type === "slot"}
            <RumusSlot value={part.value} categories={part.categories} />
        {:else}{part.value}{/if}
    {/each}
</span>

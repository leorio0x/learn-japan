<script>
    import FuriganaText from "./FuriganaText.svelte";

    let { items = [], compact = false } = $props();
</script>

{#if items.length}
    <ul class={compact ? "space-y-2" : "space-y-3"}>
        {#each items as c}
            <li class={compact ? "text-xs" : "text-sm"}>
                <p class="font-jp font-medium text-ink">
                    {#if c.pertanyaan}
                        <FuriganaText text={c.pertanyaan} />
                    {:else}
                        —
                    {/if}
                </p>
                {#if c.arti_pertanyaan}
                    <p class="text-navy/50">
                        <FuriganaText text={c.arti_pertanyaan} />
                    </p>
                {/if}
                <p class="mt-1 font-jp font-medium text-ink">
                    {#if c.jawaban}
                        <FuriganaText text={c.jawaban} />
                    {:else}
                        —
                    {/if}
                </p>
                {#if c.arti_jawaban}
                    <p class="text-navy/50">
                        <FuriganaText text={c.arti_jawaban} />
                    </p>
                {/if}
            </li>
        {/each}
    </ul>
{:else}
    <span class="text-navy/30">—</span>
{/if}

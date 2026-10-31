<script>
    let { page = 1, totalPages = 1, onPageChange = () => {} } = $props();

    function pageWindow(current, total) {
        const out = [];
        if (total <= 7) {
            for (let i = 1; i <= total; i++) out.push(i);
            return out;
        }
        out.push(1);
        const start = Math.max(2, current - 1);
        const end = Math.min(total - 1, current + 1);
        if (start > 2) out.push("…");
        for (let i = start; i <= end; i++) out.push(i);
        if (end < total - 1) out.push("…");
        out.push(total);
        return out;
    }

    const pages = $derived(pageWindow(page, totalPages));
</script>

{#if totalPages > 1}
    <div class="mt-4 flex flex-wrap items-center justify-center gap-1.5">
        <button
            class="rounded-lg border border-steel/25 bg-white px-3 py-1.5 text-xs font-semibold text-navy transition hover:bg-cream/70 disabled:cursor-not-allowed disabled:opacity-40"
            onclick={() => onPageChange(page - 1)}
            disabled={page <= 1}
        >
            ‹ Sebelumnya
        </button>

        {#each pages as p, i}
            {#if p === "…"}
                <span class="px-1 text-xs text-navy/40">…</span>
            {:else}
                <button
                    class="min-w-[2rem] rounded-lg px-2.5 py-1.5 text-xs font-semibold transition {p ===
                    page
                        ? 'bg-steel text-cream shadow-sm'
                        : 'border border-steel/25 bg-white text-navy hover:bg-cream/70'}"
                    onclick={() => onPageChange(p)}
                >
                    {p}
                </button>
            {/if}
        {/each}

        <button
            class="rounded-lg border border-steel/25 bg-white px-3 py-1.5 text-xs font-semibold text-navy transition hover:bg-cream/70 disabled:cursor-not-allowed disabled:opacity-40"
            onclick={() => onPageChange(page + 1)}
            disabled={page >= totalPages}
        >
            Berikutnya ›
        </button>
    </div>
{/if}

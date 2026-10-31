<script>
    import Furigana from "./Furigana.svelte";
    import Percakapan from "./Percakapan.svelte";

    let { item = null, onClose = () => {} } = $props();

    $effect(() => {
        function onKey(event) {
            if (event.key === "Escape") onClose();
        }
        window.addEventListener("keydown", onKey);
        return () => window.removeEventListener("keydown", onKey);
    });
</script>

{#if item}
    <div
        class="fixed inset-0 z-40 flex items-center justify-center bg-ink/60 p-4 backdrop-blur-sm"
        role="presentation"
        onclick={(e) => e.target === e.currentTarget && onClose()}
    >
        <div
            class="card max-h-[90vh] w-full max-w-md overflow-y-auto p-5 sm:p-6"
            role="dialog"
            aria-modal="true"
            aria-labelledby="detail-title"
        >
            <div class="flex items-start justify-between gap-3">
                <div class="min-w-0">
                    <h3 id="detail-title" class="text-lg font-bold text-ink">
                        Detail
                    </h3>
                    <p class="mt-0.5 text-xs text-navy/50">Seluruh data</p>
                </div>
                <button
                    type="button"
                    class="rounded-lg px-2 py-1 text-lg leading-none text-navy/50 transition hover:bg-navy/5 hover:text-navy"
                    onclick={onClose}
                    aria-label="Tutup detail"
                >
                    &times;
                </button>
            </div>

            <div class="mt-4 flex flex-wrap items-baseline gap-x-3 gap-y-1">
                <Furigana
                    kanji={item.kanji}
                    kana={item.furigana}
                    className="font-jp text-2xl font-bold text-ink"
                />
            </div>

            <div class="mt-2 flex flex-wrap gap-1.5">
                {#if item.buku}
                    <span
                        class="rounded-md bg-navy/5 px-2 py-0.5 text-[11px] font-medium text-navy/70"
                        >{item.buku}</span
                    >
                {/if}
                <span
                    class="rounded-md bg-sand/40 px-2 py-0.5 text-[11px] font-medium text-navy/70"
                    >Lesson {item.lesson}</span
                >
            </div>

            <dl class="mt-4 space-y-3 border-t border-navy/10 pt-4 text-sm">
                <div>
                    <dt class="text-xs font-semibold uppercase tracking-wide text-navy/50">
                        Kanji
                    </dt>
                    <dd class="mt-0.5 font-jp text-lg font-bold text-ink">
                        {item.kanji || "—"}
                    </dd>
                </div>
                <div>
                    <dt class="text-xs font-semibold uppercase tracking-wide text-navy/50">
                        Furigana
                    </dt>
                    <dd class="mt-0.5 font-jp text-navy/80">
                        {item.furigana || "—"}
                    </dd>
                </div>
                <div>
                    <dt class="text-xs font-semibold uppercase tracking-wide text-navy/50">
                        Arti
                    </dt>
                    <dd class="mt-0.5 font-semibold text-sage">
                        {#if (item.terjemahan || []).length}
                            <ul class="list-disc space-y-0.5 pl-4">
                                {#each item.terjemahan as arti}
                                    <li>{arti}</li>
                                {/each}
                            </ul>
                        {:else}
                            —
                        {/if}
                    </dd>
                </div>
                <div>
                    <dt class="text-xs font-semibold uppercase tracking-wide text-navy/50">
                        Contoh Percakapan
                    </dt>
                    <dd
                        class="mt-1 rounded-lg border border-navy/5 bg-cream/60 p-3"
                    >
                        <Percakapan items={item.contoh_percakapan || []} />
                    </dd>
                </div>
            </dl>

            <div class="mt-5 flex justify-end">
                <button class="btn-ghost" onclick={onClose}>Tutup</button>
            </div>
        </div>
    </div>
{/if}

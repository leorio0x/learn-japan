<script>
    import { onMount } from "svelte";
    import { api } from "../lib/api.js";
    import { parseParagraph } from "../lib/paragraph.js";
    import DokkaiWord from "../components/DokkaiWord.svelte";
    import FeatureLinks from "../components/FeatureLinks.svelte";

    let { onNavigate = () => {} } = $props();

    let levels = $state([]);
    let level = $state("");
    let stories = $state([]);
    let loading = $state(true);
    let loadingStories = $state(false);
    let error = $state("");

    onMount(loadLevels);

    async function loadLevels() {
        loading = true;
        error = "";
        try {
            levels = (await api.listDokkaiLevels()) || [];
            if (!level || !levels.includes(level)) {
                level = levels.length ? levels[0] : "";
            }
            await fetchStories();
        } catch (e) {
            error = e.message;
        } finally {
            loading = false;
        }
    }

    async function fetchStories() {
        if (!level) {
            stories = [];
            return;
        }
        loadingStories = true;
        error = "";
        try {
            stories = (await api.listDokkaiLevel(level)) || [];
        } catch (e) {
            error = e.message;
        } finally {
            loadingStories = false;
        }
    }

    async function selectLevel(event) {
        level = event.target.value;
        await fetchStories();
    }
</script>

<div class="mx-auto max-w-3xl">
    <div
        class="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"
    >
        <div>
            <button
                type="button"
                class="mb-2 inline-flex items-center gap-1 text-xs font-semibold text-steel transition hover:text-navy"
                onclick={() => onNavigate("kotoba")}
            >
                &larr; Kembali ke Kotoba
            </button>
            <h1 class="text-xl font-bold text-ink sm:text-2xl">Dokkai</h1>
            <p class="mt-1 text-sm text-navy/60">
                Bacaan bahasa Jepang berdasarkan level JLPT
            </p>
        </div>
        <FeatureLinks page="dokkai" {onNavigate} />
    </div>

    <div
        class="card mb-4 flex flex-col gap-4 p-4 sm:flex-row sm:items-end sm:p-5"
    >
        <div class="flex-1">
            <label class="label" for="dokkai-level">Pilih level</label>
            <select
                id="dokkai-level"
                class="input"
                bind:value={level}
                onchange={selectLevel}
                disabled={!levels.length}
            >
                {#if !levels.length}
                    <option value="">Tidak ada level</option>
                {/if}
                {#each levels as l}
                    <option value={l}>{l}</option>
                {/each}
            </select>
        </div>
        {#if level}
            <p class="text-xs text-navy/50">
                {stories.length} dokkai pada level {level}
            </p>
        {/if}
    </div>

    <p class="mb-3 text-xs text-navy/50">
        Klik dua kali pada kata untuk melihat furigana dan artinya.
    </p>

    {#if loading || loadingStories}
        <div class="card p-12 text-center text-navy/60">
            <span class="animate-pulse">Memuat dokkai...</span>
        </div>
    {:else if error}
        <div class="card border-red-300 p-8 text-center text-red-700">
            <p class="font-semibold">{error}</p>
            <button class="btn-primary mt-4" onclick={loadLevels}>
                Coba lagi
            </button>
        </div>
    {:else if !levels.length}
        <div class="card p-10 text-center">
            <p class="font-semibold text-ink">Belum ada level dokkai</p>
            <p class="mt-1 text-sm text-navy/60">
                Tambahkan dokkai melalui halaman Kelola Dokkai.
            </p>
        </div>
    {:else if !stories.length}
        <div class="card p-10 text-center">
            <p class="font-semibold text-ink">Belum ada dokkai</p>
            <p class="mt-1 text-sm text-navy/60">
                Belum ada dokkai untuk level {level}.
            </p>
        </div>
    {:else}
        <div class="space-y-4">
            {#each stories as story, i (story.id)}
                {@const parts = parseParagraph(story.paragraf)}
                <article class="card p-4 sm:p-6">
                    <div class="mb-3 flex flex-wrap items-center gap-2">
                        <span
                            class="rounded-md bg-sand/40 px-2 py-0.5 text-[11px] font-medium text-navy/70"
                            >{story.level || level}</span
                        >
                        {#if story.judul}
                            <h2 class="font-jp text-lg font-bold text-ink">
                                {story.judul}
                            </h2>
                        {:else}
                            <h2 class="text-lg font-bold text-ink">
                                Dokkai {i + 1}
                            </h2>
                        {/if}
                    </div>

                    <p
                        class="whitespace-pre-line font-jp text-lg leading-relaxed text-ink/90 sm:text-xl"
                    >
                        {#each parts as part}
                            {#if part.type === "text"}{part.value}{:else}<DokkaiWord
                                    kanji={part.kanji}
                                    kana={part.kana}
                                    arti={part.arti}
                                />{/if}
                        {/each}
                    </p>

                    {#if story.terjemahan}
                        <div
                            class="mt-4 rounded-xl border border-navy/5 bg-cream/60 p-3"
                        >
                            <p
                                class="mb-1 text-[11px] font-semibold uppercase tracking-wide text-navy/50"
                            >
                                Terjemahan
                            </p>
                            <p class="text-sm leading-relaxed text-navy/70">
                                {story.terjemahan}
                            </p>
                        </div>
                    {/if}
                </article>
            {/each}
        </div>
    {/if}
</div>

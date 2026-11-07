<script>
    import { onMount } from "svelte";
    import { api } from "../lib/api.js";
    import { parseParagraph } from "../lib/paragraph.js";
    import { auth } from "../lib/auth.svelte.js";
    import DokkaiWord from "./DokkaiWord.svelte";

    let { type = "kotoba", onNavigate = () => {} } = $props();

    const isKotoba = $derived(type === "kotoba");
    const title = $derived(isKotoba ? "Dokkai Kotoba" : "Dokkai Kanji");
    const subtitle = $derived(
        isKotoba
            ? "Cerita pendek yang memuat kotoba per lesson"
            : "Cerita pendek yang memuat kanji per lesson",
    );

    let lessons = $state([]);
    let stories = $state([]);
    let lesson = $state("");
    let loading = $state(true);
    let loadingStories = $state(false);
    let error = $state("");
    let hideKanji = $state(false);

    let showBulk = $state(false);
    let bulkInputEl = $state(null);
    let bulkFile = $state(null);
    let bulkUploading = $state(false);
    let bulkResult = $state(null);
    let bulkError = $state("");

    onMount(load);

    async function load() {
        loading = true;
        error = "";
        try {
            lessons = (await api.listDokkaiLessons(type)) || [];
            lesson = lessons.length ? lessons[0] : "";
            await fetchStories();
        } catch (e) {
            error = e.message;
        } finally {
            loading = false;
        }
    }

    async function fetchStories() {
        if (!lesson) {
            stories = [];
            return;
        }
        loadingStories = true;
        error = "";
        try {
            stories = (await api.getDokkai(type, lesson)) || [];
        } catch (e) {
            error = e.message;
        } finally {
            loadingStories = false;
        }
    }

    async function selectLesson(event) {
        lesson = event.target.value;
        await fetchStories();
    }

    function openBulk() {
        bulkFile = null;
        bulkResult = null;
        bulkError = "";
        if (bulkInputEl) bulkInputEl.value = "";
        showBulk = true;
    }

    function closeBulk() {
        showBulk = false;
        bulkFile = null;
        bulkResult = null;
        bulkError = "";
        if (bulkInputEl) bulkInputEl.value = "";
    }

    function onBulkFile(event) {
        bulkFile = event.target.files?.[0] || null;
        bulkError = "";
        bulkResult = null;
    }

    async function submitBulk() {
        if (!bulkFile) {
            bulkError = "Pilih file terlebih dahulu";
            return;
        }
        bulkUploading = true;
        bulkError = "";
        bulkResult = null;
        try {
            bulkResult = await api.bulkDokkai(type, bulkFile);
            bulkFile = null;
            if (bulkInputEl) bulkInputEl.value = "";
            await fetchStories();
        } catch (e) {
            bulkError = e.message;
        } finally {
            bulkUploading = false;
        }
    }
</script>

<div class="mx-auto max-w-3xl">
    <div class="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
            <button
                type="button"
                class="mb-2 inline-flex items-center gap-1 text-xs font-semibold text-steel transition hover:text-navy"
                onclick={() => onNavigate(isKotoba ? "kotoba" : "kanji")}
            >
                ‹ Kembali ke {isKotoba ? "Kotoba" : "Kanji"}
            </button>
            <h1 class="text-xl font-bold text-ink sm:text-2xl">{title}</h1>
            <p class="mt-1 text-sm text-navy/60">{subtitle}</p>
        </div>

        {#if auth.authenticated}
            <button
                type="button"
                class="btn-ghost w-full sm:w-auto"
                onclick={openBulk}
            >
                Bulk Input Dokkai
            </button>
        {/if}
    </div>

    <div class="card mb-4 flex flex-col gap-4 p-4 sm:flex-row sm:items-end sm:p-5">
        <div class="flex-1">
            <label class="label" for="dokkai-lesson">Filter lesson</label>
            <select
                id="dokkai-lesson"
                class="input"
                bind:value={lesson}
                onchange={selectLesson}
            >
                {#if !lessons.length}
                    <option value="">Tidak ada lesson</option>
                {/if}
                {#each lessons as l}
                    <option value={l}>Lesson {l}</option>
                {/each}
            </select>
        </div>

        {#if isKotoba}
            <button
                type="button"
                class="flex items-center gap-2.5 rounded-xl border border-steel/30 bg-white px-3.5 py-2.5 text-sm font-semibold text-navy transition hover:bg-cream/60"
                onclick={() => (hideKanji = !hideKanji)}
                aria-pressed={hideKanji}
            >
                <span
                    class="relative inline-flex h-5 w-9 items-center rounded-full transition {hideKanji
                        ? 'bg-sage'
                        : 'bg-navy/20'}"
                >
                    <span
                        class="absolute h-4 w-4 rounded-full bg-white shadow transition-all {hideKanji
                            ? 'left-4'
                            : 'left-0.5'}"
                    ></span>
                </span>
                Sembunyikan Kanji
            </button>
        {/if}
    </div>

    <p class="mb-3 text-xs text-navy/50">
        Klik pada kata untuk melihat furigana dan artinya.
    </p>

    {#if loading || loadingStories}
        <div class="card p-12 text-center text-navy/60">
            <span class="animate-pulse">Memuat dokkai...</span>
        </div>
    {:else if error}
        <div class="card border-red-300 p-8 text-center text-red-700">
            <p class="font-semibold">{error}</p>
            <button class="btn-primary mt-4" onclick={load}>Coba lagi</button>
        </div>
    {:else if !stories.length}
        <div class="card p-10 text-center">
            <p class="font-semibold text-ink">Belum ada dokkai</p>
            <p class="mt-1 text-sm text-navy/60">
                Belum ada dokkai untuk lesson ini.
            </p>
        </div>
    {:else}
        <div class="space-y-4">
            {#each stories as story, i}
                {@const parts = parseParagraph(story.paragraf)}
                <article class="card p-4 sm:p-6">
                    <div class="mb-3 flex flex-wrap items-center gap-2">
                        <span
                            class="rounded-md bg-sand/40 px-2 py-0.5 text-[11px] font-medium text-navy/70"
                            >Lesson {story.lesson || lesson}</span
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
                                    {hideKanji}
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
                            <p
                                class="whitespace-pre-line text-sm leading-relaxed text-navy/70"
                            >
                                {story.terjemahan}
                            </p>
                        </div>
                    {/if}
                </article>
            {/each}
        </div>
    {/if}
</div>

{#if showBulk}
    <div
        class="fixed inset-0 z-40 flex items-center justify-center bg-ink/60 p-4 backdrop-blur-sm"
    >
        <div class="card w-full max-w-md p-5 sm:p-6">
            <h3 class="text-lg font-bold text-ink">Bulk Input Dokkai</h3>
            <p class="mt-2 text-sm text-navy/70">
                Unggah file <span class="font-semibold">.json</span> berisi array
                objek sesuai format
                <span class="font-mono">data/format/dokkai.json</span>.
            </p>
            <p
                class="mt-1.5 overflow-x-auto rounded-lg bg-navy/5 px-3 py-2 font-mono text-xs text-navy/80"
            >
                [&#123;"lesson":"1","judul":"...","paragraf":"&lt;kanji:kana:arti&gt;","terjemahan":"..."&#125;]
            </p>
            <p class="mt-2 text-xs text-navy/50">
                Data akan ditambahkan ke
                <span class="font-mono">data/dokkai/{type}/[lesson].json</span>
                sesuai nilai <span class="font-mono">lesson</span>.
            </p>

            <input
                type="file"
                accept=".json,application/json"
                bind:this={bulkInputEl}
                onchange={onBulkFile}
                class="mt-4 block w-full text-sm text-navy file:mr-3 file:rounded-lg file:border-0 file:bg-steel file:px-3.5 file:py-2 file:text-sm file:font-semibold file:text-cream hover:file:bg-navy"
            />

            {#if bulkError}
                <p
                    class="mt-3 rounded-lg bg-red-500/10 px-3 py-2 text-sm text-red-700"
                >
                    {bulkError}
                </p>
            {/if}
            {#if bulkResult}
                <p
                    class="mt-3 rounded-lg bg-sage/15 px-3 py-2 text-sm font-medium text-navy"
                >
                    Berhasil mengimpor {bulkResult.imported} dokkai.
                </p>
            {/if}

            <div
                class="mt-5 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"
            >
                <button
                    class="btn-ghost"
                    onclick={closeBulk}
                    disabled={bulkUploading}>Tutup</button
                >
                <button
                    class="btn-primary"
                    onclick={submitBulk}
                    disabled={bulkUploading || !bulkFile}
                >
                    {bulkUploading ? "Mengunggah..." : "Unggah"}
                </button>
            </div>
        </div>
    </div>
{/if}

<script>
    import { onMount } from "svelte";
    import { api } from "../lib/api.js";
    import { uniqueSortedLessons } from "../lib/quiz.js";
    import Furigana from "../components/Furigana.svelte";

    let items = $state([]);
    let loading = $state(true);
    let error = $state("");
    let lesson = $state("1");
    let hideTerjemahan = $state(false);

    const lessons = $derived(uniqueSortedLessons(items));
    const filtered = $derived(
        lesson === "all"
            ? items
            : items.filter((k) => String(k.lesson) === lesson),
    );

    onMount(load);

    async function load() {
        loading = true;
        error = "";
        try {
            items = (await api.listKaiwa()) || [];
            const available = uniqueSortedLessons(items);
            if (lesson !== "all" && !available.includes(lesson)) {
                lesson = available.length ? available[0] : "all";
            }
        } catch (e) {
            error = e.message;
        } finally {
            loading = false;
        }
    }
</script>

<div class="mx-auto max-w-3xl">
    <div class="mb-5">
        <h1 class="text-xl font-bold text-ink sm:text-2xl">Daftar Kaiwa</h1>
        <p class="mt-1 text-sm text-navy/60">
            Lihat daftar percakapan bahasa Jepang
        </p>
    </div>

    <div
        class="card mb-4 flex flex-col gap-4 p-4 sm:flex-row sm:items-end sm:p-5"
    >
        <div class="flex-1">
            <label class="label" for="lesson">Filter lesson</label>
            <select id="lesson" class="input" bind:value={lesson}>
                <option value="all">Semua lesson</option>
                {#each lessons as l}
                    <option value={l}>Lesson {l}</option>
                {/each}
            </select>
        </div>

        <button
            type="button"
            class="flex items-center gap-2.5 rounded-xl border border-steel/30 bg-white px-3.5 py-2.5 text-sm font-semibold text-navy transition hover:bg-cream/60"
            onclick={() => (hideTerjemahan = !hideTerjemahan)}
            aria-pressed={hideTerjemahan}
        >
            <span
                class="relative inline-flex h-5 w-9 items-center rounded-full transition {hideTerjemahan
                    ? 'bg-sage'
                    : 'bg-navy/20'}"
            >
                <span
                    class="absolute h-4 w-4 rounded-full bg-white shadow transition-all {hideTerjemahan
                        ? 'left-4'
                        : 'left-0.5'}"
                ></span>
            </span>
            Sembunyikan Terjemahan
        </button>
    </div>

    {#if loading}
        <div class="card p-12 text-center text-navy/60">
            <span class="animate-pulse">Memuat data...</span>
        </div>
    {:else if error}
        <div class="card border-red-300 p-8 text-center text-red-700">
            <p class="font-semibold">{error}</p>
            <button class="btn-primary mt-4" onclick={load}>Coba lagi</button>
        </div>
    {:else if !filtered.length}
        <div class="card p-10 text-center">
            <p class="font-semibold text-ink">Tidak ada kaiwa</p>
            <p class="mt-1 text-sm text-navy/60">
                Belum ada kaiwa untuk lesson ini.
            </p>
        </div>
    {:else}
        <ul class="space-y-3">
            {#each filtered as k (k.id)}
                <li class="card p-4 sm:p-5">
                    <div class="mb-3 flex flex-wrap gap-1.5">
                        {#if k.buku}
                            <span
                                class="rounded-md bg-navy/5 px-2 py-0.5 text-[11px] font-medium text-navy/70"
                                >{k.buku}</span
                            >
                        {/if}
                        <span
                            class="rounded-md bg-sand/40 px-2 py-0.5 text-[11px] font-medium text-navy/70"
                            >Lesson {k.lesson}</span
                        >
                    </div>

                    <div class="space-y-4">
                        {#each k.pola || [] as p}
                            <div>
                                {#if p.pola_pertanyaan ||
                                    p.pola_jawaban ||
                                    p.deskripsi}
                                    <div class="mb-2">
                                        <p
                                            class="text-xs font-semibold text-navy/70"
                                        >
                                            <span class="font-jp text-ink"
                                                >{p.pola_pertanyaan}</span
                                            >
                                            <span class="text-navy/40">→</span>
                                            <span class="font-jp text-ink"
                                                >{p.pola_jawaban}</span
                                            >
                                        </p>
                                        {#if p.deskripsi}
                                            <p class="text-xs text-navy/50">
                                                {p.deskripsi}
                                            </p>
                                        {/if}
                                    </div>
                                {/if}

                                {#each p.contoh || [] as c}
                                    <div
                                        class="space-y-2.5 rounded-lg border border-navy/10 bg-cream/30 p-3"
                                    >
                                        {#each c.pertanyaan || [] as q}
                                            <div
                                                class="flex items-start gap-2.5"
                                            >
                                                <span
                                                    class="mt-0.5 shrink-0 rounded-md bg-steel/10 px-1.5 py-0.5 text-[10px] font-bold uppercase text-navy/60"
                                                    >P :</span
                                                >
                                                <div class="min-w-0">
                                                    <Furigana
                                                        kanji={q.kanji}
                                                        kana={q[
                                                            "hiragana/katakana"
                                                        ]}
                                                        className="font-jp text-lg font-semibold text-ink"
                                                    />
                                                    {#if !hideTerjemahan &&
                                                        q.terjemahan}
                                                        <p
                                                            class="text-sm text-navy/60"
                                                        >
                                                            {q.terjemahan}
                                                        </p>
                                                    {/if}
                                                </div>
                                            </div>
                                        {/each}

                                        {#each c.jawaban || [] as j}
                                            <div
                                                class="flex items-start gap-2.5"
                                            >
                                                <span
                                                    class="mt-0.5 shrink-0 rounded-md bg-sage/15 px-1.5 py-0.5 text-[10px] font-bold uppercase text-navy/60"
                                                    >J :</span
                                                >
                                                <div class="min-w-0">
                                                    <Furigana
                                                        kanji={j.kanji}
                                                        kana={j[
                                                            "hiragana/katakana"
                                                        ]}
                                                        className="font-jp text-lg font-semibold text-sage"
                                                    />
                                                    {#if !hideTerjemahan &&
                                                        j.terjemahan}
                                                        <p
                                                            class="text-sm text-navy/60"
                                                        >
                                                            {j.terjemahan}
                                                        </p>
                                                    {/if}
                                                </div>
                                            </div>
                                        {/each}
                                    </div>
                                {/each}
                            </div>
                        {/each}
                    </div>
                </li>
            {/each}
        </ul>
    {/if}
</div>

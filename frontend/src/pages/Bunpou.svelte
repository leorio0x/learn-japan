<script>
    import { onMount } from "svelte";
    import { api } from "../lib/api.js";
    import { lessonLabel, sortLessons } from "../lib/bunpou.js";
    import BunpouRumus from "../components/BunpouRumus.svelte";
    import BunpouSentence from "../components/BunpouSentence.svelte";
    import FeatureLinks from "../components/FeatureLinks.svelte";
    import Pagination from "../components/Pagination.svelte";

    let { onNavigate = () => {} } = $props();

    let items = $state([]);
    let loading = $state(true);
    let error = $state("");
    let lesson = $state("all");
    let hideKanji = $state(false);

    const pageSize = 5;
    let page = $state(1);

    const lessons = $derived(sortLessons(items));

    const filtered = $derived(
        lesson === "all"
            ? items
            : items.filter((b) => String(b.lesson) === lesson),
    );

    const patterns = $derived(
        filtered.flatMap((b, bi) =>
            (b.pola || []).map((pola, index) => ({
                key: `${b.id ?? `bunpou-${bi}`}-${index}`,
                bunpou: b,
                pola,
            })),
        ),
    );

    const totalPages = $derived(Math.max(1, Math.ceil(patterns.length / pageSize)));
    const currentPage = $derived(Math.min(page, totalPages));
    const paged = $derived(
        patterns.slice((currentPage - 1) * pageSize, currentPage * pageSize),
    );
    const rangeStart = $derived((currentPage - 1) * pageSize + 1);
    const rangeEnd = $derived(
        Math.min(currentPage * pageSize, patterns.length),
    );

    onMount(load);

    async function load() {
        loading = true;
        error = "";
        try {
            items = (await api.listBunpou()) || [];
            if (
                lesson !== "all" &&
                !items.some((b) => String(b.lesson) === lesson)
            ) {
                lesson = "all";
            }
        } catch (e) {
            error = e.message;
        } finally {
            loading = false;
        }
    }
</script>

<div class="mx-auto max-w-3xl">
    <div
        class="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"
    >
        <div>
            <h1 class="text-xl font-bold text-ink sm:text-2xl">Daftar Bunpou</h1>
            <p class="mt-1 text-sm text-navy/60">
                Pelajari pola kalimat bahasa Jepang
            </p>
        </div>
        <FeatureLinks page="bunpou" {onNavigate} />
    </div>

    <div
        class="card mb-4 flex flex-col gap-4 p-4 sm:flex-row sm:items-end sm:p-5"
    >
        <div class="flex-1">
            <label class="label" for="lesson">Filter lesson</label>
            <select
                id="lesson"
                class="input"
                bind:value={lesson}
                onchange={() => (page = 1)}
            >
                <option value="all">Semua lesson</option>
                {#each lessons as l}
                    <option value={l}>{lessonLabel(l)}</option>
                {/each}
            </select>
        </div>

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
    </div>

    <p class="mb-3 text-xs text-navy/50">
        Klik dua kali pada kata untuk melihat furigana dan arti. Klik simbol ～
        pada rumus untuk melihat kategori pengisi.
    </p>

    {#if loading}
        <div class="card p-12 text-center text-navy/60">
            <span class="animate-pulse">Memuat data...</span>
        </div>
    {:else if error}
        <div class="card border-red-300 p-8 text-center text-red-700">
            <p class="font-semibold">{error}</p>
            <button class="btn-primary mt-4" onclick={load}>Coba lagi</button>
        </div>
    {:else if !patterns.length}
        <div class="card p-10 text-center">
            <p class="font-semibold text-ink">Tidak ada bunpou</p>
            <p class="mt-1 text-sm text-navy/60">
                Belum ada bunpou untuk lesson ini.
            </p>
        </div>
    {:else}
        <ul class="space-y-3">
            {#each paged as item (item.key)}
                {@const b = item.bunpou}
                {@const pola = item.pola}
                <li class="card p-4 sm:p-5">
                    <div class="mb-3 flex flex-wrap gap-1.5">
                        {#if b.buku}
                            <span
                                class="rounded-md bg-navy/5 px-2 py-0.5 text-[11px] font-medium text-navy/70"
                                >{b.buku}</span
                            >
                        {/if}
                        <span
                            class="rounded-md bg-sand/40 px-2 py-0.5 text-[11px] font-medium text-navy/70"
                            >{lessonLabel(b.lesson)}</span
                        >
                        {#if pola.no_bunpou}
                            <span
                                class="rounded-md bg-sage/15 px-2 py-0.5 text-[11px] font-semibold text-navy/70"
                                >No. {pola.no_bunpou}</span
                            >
                        {/if}
                    </div>

                    {#if pola.nama_pola}
                        <h2 class="text-lg font-bold text-ink">
                            {pola.nama_pola}
                        </h2>
                    {/if}

                    {#if pola.rumus}
                        <div
                            class="mt-2 rounded-xl border border-steel/20 bg-cream/50 px-3.5 py-2.5"
                        >
                            <p
                                class="text-[10px] font-semibold uppercase tracking-wide text-navy/45"
                            >
                                Rumus
                            </p>
                            <BunpouRumus
                                rumus={pola.rumus}
                                className="text-base text-ink sm:text-lg"
                            />
                        </div>
                    {/if}

                    {#if pola.deskripsi}
                        <p class="mt-2 text-sm leading-relaxed text-navy/65">
                            {pola.deskripsi}
                        </p>
                    {/if}

                    {#if pola.contoh?.length}
                        <div class="mt-3 space-y-2.5">
                            {#each pola.contoh as contoh, ci}
                                <div
                                    class="rounded-lg border border-navy/10 bg-cream/30 p-3"
                                >
                                    {#if contoh.pertanyaan}
                                        <div class="flex items-start gap-2.5">
                                            <span
                                                class="mt-1 shrink-0 rounded-md bg-steel/10 px-1.5 py-0.5 text-[10px] font-bold uppercase text-navy/60"
                                                >P</span
                                            >
                                            <BunpouSentence
                                                kalimat={contoh.pertanyaan.kalimat}
                                                terjemahan={contoh.pertanyaan
                                                    .terjemahan}
                                                {hideKanji}
                                                className="text-lg font-semibold text-ink"
                                            />
                                        </div>
                                    {/if}

                                    {#if contoh.jawaban}
                                        <div
                                            class="flex items-start gap-2.5 {contoh.pertanyaan
                                                ? 'mt-2.5'
                                                : ''}"
                                        >
                                            <span
                                                class="mt-1 shrink-0 rounded-md bg-sage/15 px-1.5 py-0.5 text-[10px] font-bold uppercase text-navy/60"
                                                >{contoh.pertanyaan ? "J" : "Contoh"}</span
                                            >
                                            <BunpouSentence
                                                kalimat={contoh.jawaban.kalimat}
                                                terjemahan={contoh.jawaban
                                                    .terjemahan}
                                                {hideKanji}
                                                className="text-lg font-semibold text-sage"
                                            />
                                        </div>
                                    {/if}
                                </div>
                            {/each}
                        </div>
                    {/if}
                </li>
            {/each}
        </ul>

        <p class="mt-3 text-center text-xs text-navy/50">
            Menampilkan {rangeStart}-{rangeEnd} dari {patterns.length} pola
        </p>
        <Pagination
            page={currentPage}
            {totalPages}
            onPageChange={(p) => (page = p)}
        />
    {/if}
</div>

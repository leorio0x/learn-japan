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
        class="mb-4 flex flex-col gap-4 rounded-2xl border border-steel/20 bg-cream/80 p-4 shadow-sm backdrop-blur sm:flex-row sm:items-end sm:p-5"
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
    </div>

    <p class="mb-3 text-xs text-navy/50">
        Klik pada kata untuk melihat furigana dan arti. Klik simbol ～
        pada rumus untuk melihat kategori pengisi.
    </p>

    {#if loading}
        <div
            class="rounded-2xl border border-steel/20 bg-cream/80 p-12 text-center text-navy/60 shadow-sm"
        >
            <span class="animate-pulse">Memuat data...</span>
        </div>
    {:else if error}
        <div
            class="rounded-2xl border border-sand bg-cream/80 p-8 text-center text-ink shadow-sm"
        >
            <p class="font-semibold">{error}</p>
            <button class="btn-primary mt-4" onclick={load}>Coba lagi</button>
        </div>
    {:else if !patterns.length}
        <div
            class="rounded-2xl border border-steel/20 bg-cream/80 p-10 text-center shadow-sm"
        >
            <p class="font-semibold text-ink">Tidak ada bunpou</p>
            <p class="mt-1 text-sm text-navy/60">
                Belum ada bunpou untuk lesson ini.
            </p>
        </div>
    {:else}
        <ul class="space-y-4">
            {#each paged as item, i (item.key)}
                {@const b = item.bunpou}
                {@const pola = item.pola}
                {@const nomor = pola.no_bunpou || String(rangeStart + i).padStart(2, "0")}
                <li
                    class="group relative overflow-hidden rounded-3xl border border-steel/20 bg-cream shadow-sm transition duration-300 hover:-translate-y-0.5 hover:border-sage/50 hover:shadow-md"
                >
                    <div class="grid sm:grid-cols-[10rem_1fr]">
                        <div
                            class="relative overflow-hidden bg-gradient-to-br from-ink via-navy to-steel p-4 text-cream sm:p-5"
                        >
                            <span
                                class="pointer-events-none absolute -right-4 -top-5 select-none font-jp text-7xl font-black leading-none text-cream/10"
                                >文</span
                            >

                            <div
                                class="relative flex items-center gap-3 sm:flex-col sm:items-start sm:gap-0"
                            >
                                <span
                                    class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-sand text-lg font-black text-ink shadow-md"
                                    >{nomor}</span
                                >
                                <div class="min-w-0">
                                    <p
                                        class="text-[9px] font-bold uppercase tracking-[0.25em] text-sand/80"
                                    >
                                        Bunpou
                                    </p>
                                    <p
                                        class="mt-0.5 truncate text-sm font-semibold text-cream"
                                    >
                                        {lessonLabel(b.lesson)}
                                    </p>
                                    {#if b.buku}
                                        <p
                                            class="mt-0.5 text-[10px] leading-snug text-cream/50"
                                        >
                                            {b.buku}
                                        </p>
                                    {/if}
                                </div>
                            </div>

                            {#if pola.nama_pola}
                                <div
                                    class="relative mt-4 hidden border-t border-dashed border-cream/15 pt-3 sm:block"
                                >
                                    <p
                                        class="text-[9px] font-bold uppercase tracking-[0.2em] text-cream/40"
                                    >
                                        Pola
                                    </p>
                                    <p
                                        class="mt-1 text-sm leading-snug text-cream/80"
                                    >
                                        {pola.nama_pola}
                                    </p>
                                </div>
                            {/if}
                        </div>

                        <div class="relative p-4 sm:p-5">
                            {#if pola.nama_pola}
                                <h2
                                    class="text-lg font-bold text-ink sm:hidden"
                                >
                                    {pola.nama_pola}
                                </h2>
                            {/if}

                            {#if pola.rumus}
                                <div
                                    class="mt-3 rounded-xl border border-dashed border-steel/30 bg-sand/30 px-3.5 py-2.5"
                                >
                                    <p
                                        class="text-[9px] font-bold uppercase tracking-[0.2em] text-steel/70"
                                    >
                                        Rumus
                                    </p>
                                    <BunpouRumus
                                        rumus={pola.rumus}
                                        className="mt-1 block text-base leading-relaxed text-ink sm:text-lg"
                                    />
                                </div>
                            {/if}

                            {#if pola.deskripsi}
                                <p
                                    class="mt-3 text-sm leading-relaxed text-navy/65"
                                >
                                    {pola.deskripsi}
                                </p>
                            {/if}

                            {#if pola.contoh?.length}
                                <div class="mt-4 space-y-2.5">
                                    {#each pola.contoh as contoh, ci}
                                        <div
                                            class="rounded-xl border border-steel/15 bg-sand/20 p-3"
                                        >
                                            <div
                                                class="mb-2 flex items-center gap-2"
                                            >
                                                <span
                                                    class="rounded bg-sand/60 px-1.5 py-0.5 text-[9px] font-bold uppercase tracking-wide text-navy/60"
                                                    >Contoh {ci + 1}</span
                                                >
                                                <span
                                                    class="h-px flex-1 bg-gradient-to-r from-steel/25 to-transparent"
                                                ></span>
                                            </div>

                                            <div class="space-y-2">
                                                {#if contoh.pertanyaan}
                                                    <div
                                                        class="flex items-start gap-2"
                                                    >
                                                        <span
                                                            class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-md bg-steel text-[10px] font-bold text-cream"
                                                            >P</span
                                                        >
                                                        <div
                                                            class="min-w-0 flex-1 rounded-lg bg-cream px-3 py-2 ring-1 ring-steel/10"
                                                        >
                                                            <BunpouSentence
                                                                kalimat={contoh
                                                                    .pertanyaan
                                                                    .kalimat}
                                                                terjemahan={contoh
                                                                    .pertanyaan
                                                                    .terjemahan}
                                                                className="text-base font-semibold text-ink sm:text-lg"
                                                            />
                                                        </div>
                                                    </div>
                                                {/if}

                                                {#if contoh.jawaban}
                                                    <div
                                                        class="flex items-start gap-2"
                                                    >
                                                        <span
                                                            class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-md bg-sage text-[10px] font-bold text-cream"
                                                            >{contoh.pertanyaan
                                                                ? "J"
                                                                : "C"}</span
                                                        >
                                                        <div
                                                            class="min-w-0 flex-1 rounded-lg bg-sage/10 px-3 py-2 ring-1 ring-sage/20"
                                                        >
                                                            <BunpouSentence
                                                                kalimat={contoh
                                                                    .jawaban
                                                                    .kalimat}
                                                                terjemahan={contoh
                                                                    .jawaban
                                                                    .terjemahan}
                                                                className="text-base font-semibold text-navy sm:text-lg"
                                                            />
                                                        </div>
                                                    </div>
                                                {/if}
                                            </div>
                                        </div>
                                    {/each}
                                </div>
                            {/if}
                        </div>
                    </div>
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

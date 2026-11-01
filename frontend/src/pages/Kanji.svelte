<script>
    import { onMount } from "svelte";
    import { api } from "../lib/api.js";
    import { uniqueSortedLessons } from "../lib/quiz.js";
    import Furigana from "../components/Furigana.svelte";
    import Pagination from "../components/Pagination.svelte";
    import Percakapan from "../components/Percakapan.svelte";
    import DetailModal from "../components/DetailModal.svelte";
    import FeatureLinks from "../components/FeatureLinks.svelte";

    let { onNavigate = () => {} } = $props();

    let items = $state([]);
    let detailItem = $state(null);
    let loading = $state(true);
    let error = $state("");
    let lesson = $state("1");
    let page = $state(1);
    let view = $state("card");
    let hideKanji = $state(false);
    let hideFurigana = $state(false);
    let hideArti = $state(false);
    let hidePercakapan = $state(false);
    let flashIndex = $state(0);
    let flipped = $state(false);

    function toggleClass(active) {
        return `flex-1 rounded-lg px-3 py-2 text-xs font-semibold transition ${
            active
                ? "bg-steel text-cream shadow-sm"
                : "text-navy/70 hover:bg-white"
        }`;
    }

    const pageSize = 10;

    const viewOptions = [
        { id: "card", label: "Card" },
        { id: "table", label: "Table" },
        { id: "flashcard", label: "Flashcard" },
    ];

    const lessons = $derived(uniqueSortedLessons(items));
    const filtered = $derived(
        lesson === "all"
            ? items
            : items.filter((k) => String(k.lesson) === lesson),
    );
    const paginate = $derived(lesson === "all");
    const totalPages = $derived(
        paginate ? Math.max(1, Math.ceil(filtered.length / pageSize)) : 1,
    );
    const currentPage = $derived(paginate ? Math.min(page, totalPages) : 1);
    const paged = $derived(
        paginate
            ? filtered.slice(
                  (currentPage - 1) * pageSize,
                  currentPage * pageSize,
              )
            : filtered,
    );
    const rangeStart = $derived((currentPage - 1) * pageSize + 1);
    const rangeEnd = $derived(
        paginate
            ? Math.min(currentPage * pageSize, filtered.length)
            : filtered.length,
    );
    const flashTotal = $derived(filtered.length);
    const flashCurrent = $derived(
        flashTotal ? filtered[Math.min(flashIndex, flashTotal - 1)] : null,
    );

    function resetFlash() {
        flashIndex = 0;
        flipped = false;
    }

    function selectLesson() {
        page = 1;
        resetFlash();
    }

    function selectView(id) {
        view = id;
        if (id === "flashcard") resetFlash();
    }

    function goFlash(delta) {
        if (!flashTotal) return;
        flashIndex =
            (Math.min(flashIndex, flashTotal - 1) + delta + flashTotal) %
            flashTotal;
        flipped = false;
    }

    onMount(load);

    async function load() {
        loading = true;
        error = "";
        try {
            items = (await api.listKanji()) || [];
            const available = uniqueSortedLessons(items);
            if (lesson !== "all" && !available.includes(lesson)) {
                lesson = available.length ? available[0] : "all";
            }
            page = 1;
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
            <h1 class="text-xl font-bold text-ink sm:text-2xl">Daftar Kanji</h1>
            <p class="mt-1 text-sm text-navy/60">Lihat daftar kanji</p>
        </div>
        <FeatureLinks page="kanji" {onNavigate} />
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
                onchange={selectLesson}
            >
                <option value="all">Semua lesson</option>
                {#each lessons as l}
                    <option value={l}>Lesson {l}</option>
                {/each}
            </select>
        </div>

        <div>
            <span class="label">Tampilan</span>
            <div class="flex gap-1 rounded-xl bg-navy/5 p-1">
                {#each viewOptions as v}
                    <button
                        class="flex-1 rounded-lg px-3.5 py-2 text-xs font-semibold transition sm:text-sm {view ===
                        v.id
                            ? 'bg-steel text-cream shadow-sm'
                            : 'text-navy/70 hover:bg-white'}"
                        onclick={() => selectView(v.id)}
                    >
                        {v.label}
                    </button>
                {/each}
            </div>
        </div>

        <div>
            <span class="label">Sembunyikan</span>
            <div class="flex gap-1 rounded-xl bg-navy/5 p-1">
                <button
                    type="button"
                    class={toggleClass(hideKanji)}
                    onclick={() => (hideKanji = !hideKanji)}
                    aria-pressed={hideKanji}
                >
                    Kanji
                </button>
                <button
                    type="button"
                    class={toggleClass(hideFurigana)}
                    onclick={() => (hideFurigana = !hideFurigana)}
                    aria-pressed={hideFurigana}
                >
                    Furigana
                </button>
                <button
                    type="button"
                    class={toggleClass(hideArti)}
                    onclick={() => (hideArti = !hideArti)}
                    aria-pressed={hideArti}
                >
                    Arti
                </button>
                <button
                    type="button"
                    class={toggleClass(hidePercakapan)}
                    onclick={() => (hidePercakapan = !hidePercakapan)}
                    aria-pressed={hidePercakapan}
                >
                    Contoh
                </button>
            </div>
        </div>
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
            <p class="font-semibold text-ink">Tidak ada kanji</p>
            <p class="mt-1 text-sm text-navy/60">
                Belum ada kanji untuk lesson ini.
            </p>
        </div>
    {:else if view === "card"}
        <ul class="space-y-2.5">
            {#each paged as k (k.id)}
                <li class="card p-3.5 sm:p-4">
                    <div class="flex items-center gap-3 sm:gap-4">
                        <div class="min-w-0 flex-1">
                            <div
                                class="flex flex-wrap items-baseline gap-x-3 gap-y-1"
                            >
                                {#if !hideKanji || !hideFurigana}
                                    <Furigana
                                        kanji={k.kanji}
                                        kana={hideFurigana ? "" : k.furigana}
                                        {hideKanji}
                                        className="font-jp text-xl font-bold text-ink"
                                    />
                                {/if}
                                {#if !hideArti}
                                    <span
                                        class="text-base font-semibold text-sage sm:text-lg"
                                        >{(k.terjemahan || []).join(", ")}</span
                                    >
                                {/if}
                            </div>
                            <div class="mt-1.5 flex flex-wrap gap-1.5">
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
                            {#if !hidePercakapan && (k.contoh_percakapan || []).length}
                                <div
                                    class="mt-2.5 rounded-lg border border-navy/5 bg-cream/60 p-2.5"
                                >
                                    <Percakapan
                                        items={k.contoh_percakapan}
                                        compact
                                    />
                                </div>
                            {/if}
                        </div>
                    </div>
                    <div class="mt-2.5">
                        <button
                            type="button"
                            class="rounded-lg bg-sage px-3 py-1.5 text-xs font-semibold text-cream shadow-sm transition hover:opacity-90"
                            onclick={() => (detailItem = k)}
                        >
                            Lihat detail
                        </button>
                    </div>
                </li>
            {/each}
        </ul>
    {:else if view === "table"}
        <div class="card overflow-hidden">
            <div class="overflow-x-auto">
                <table class="w-full min-w-[720px] text-left text-sm">
                    <thead class="bg-navy text-cream">
                        <tr>
                            {#if !hideKanji}
                                <th class="px-4 py-3 font-semibold">Kanji</th>
                            {/if}
                            {#if !hideFurigana}
                                <th class="px-4 py-3 font-semibold">
                                    Furigana
                                </th>
                            {/if}
                            {#if !hideArti}
                                <th class="px-4 py-3 font-semibold">Arti</th>
                            {/if}
                            {#if !hidePercakapan}
                                <th class="px-4 py-3 font-semibold">
                                    Contoh Percakapan
                                </th>
                            {/if}
                        </tr>
                    </thead>
                    <tbody>
                        {#each paged as k (k.id)}
                            <tr
                                class="border-t border-navy/5 odd:bg-white even:bg-cream/50"
                            >
                                {#if !hideKanji}
                                    <td
                                        class="px-4 py-3 font-jp text-lg font-bold text-ink"
                                    >
                                        {k.kanji || "—"}
                                    </td>
                                {/if}
                                {#if !hideFurigana}
                                    <td class="px-4 py-3 font-jp text-navy/80">
                                        {k.furigana}
                                    </td>
                                {/if}
                                {#if !hideArti}
                                    <td
                                        class="px-4 py-3 font-semibold text-sage"
                                    >
                                        {(k.terjemahan || []).join(", ")}
                                    </td>
                                {/if}
                                {#if !hidePercakapan}
                                    <td
                                        class="min-w-[240px] px-4 py-3 align-top"
                                    >
                                        <Percakapan
                                            items={k.contoh_percakapan}
                                            compact
                                        />
                                    </td>
                                {/if}
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        </div>
    {:else}
        <div class="mx-auto max-w-xl">
            <div
                class="mb-2 flex items-center justify-between text-sm font-semibold text-navy/70"
            >
                <span>Kartu {flashIndex + 1} / {flashTotal}</span>
                <span class="text-xs font-normal text-navy/50"
                    >Klik kartu untuk membalik</span
                >
            </div>
            <div class="h-2 w-full overflow-hidden rounded-full bg-navy/10">
                <div
                    class="h-full rounded-full bg-steel transition-all"
                    style="width: {((flashIndex + 1) / flashTotal) * 100}%"
                ></div>
            </div>

            <div
                class="flip-card mt-4 select-none"
                role="button"
                tabindex="0"
                aria-label="Balik kartu"
                onclick={() => (flipped = !flipped)}
                onkeydown={(e) => {
                    if (e.key === "Enter" || e.key === " ") {
                        e.preventDefault();
                        flipped = !flipped;
                    }
                }}
            >
                <div
                    class="flip-card-inner h-80 sm:h-96 {flipped
                        ? 'is-flipped'
                        : ''}"
                >
                    <div
                        class="flip-face card absolute inset-0 flex flex-col items-center justify-center p-6 text-center"
                    >
                        {#if hideKanji}
                            <span class="text-2xl text-navy/30">—</span>
                        {:else}
                            <span
                                class="font-jp text-7xl font-bold text-ink sm:text-8xl"
                                >{flashCurrent.kanji || flashCurrent.furigana}</span
                            >
                        {/if}
                    </div>

                    <div
                        class="flip-face flip-face-back card absolute inset-0 flex flex-col items-center justify-center gap-4 overflow-y-auto p-6 text-center"
                    >
                        {#if !hideFurigana}
                            <p class="font-jp text-3xl font-semibold text-navy/80">
                                {flashCurrent.furigana || "—"}
                            </p>
                        {/if}
                        {#if !hideArti}
                            <p
                                class="text-xl font-semibold text-sage sm:text-2xl"
                            >
                                {(flashCurrent.terjemahan || []).join(", ") ||
                                    "—"}
                            </p>
                        {/if}
                    </div>
                </div>
            </div>

            <div class="mt-4 flex items-center justify-between gap-2">
                <button
                    class="btn-ghost px-3 py-2 text-xs sm:px-4 sm:text-sm"
                    onclick={() => goFlash(-1)}
                >
                    ‹ Sebelumnya
                </button>
                <button
                    class="btn-ghost px-3 py-2 text-xs sm:px-4 sm:text-sm"
                    onclick={() => (flipped = !flipped)}
                >
                    {flipped ? "Lihat soal" : "Lihat jawaban"}
                </button>
                <button
                    class="btn-ghost px-3 py-2 text-xs sm:px-4 sm:text-sm"
                    onclick={() => goFlash(1)}
                >
                    Berikutnya ›
                </button>
            </div>

            <div class="mt-3 text-center">
                <button
                    type="button"
                    class="rounded-lg bg-sage px-3 py-1.5 text-xs font-semibold text-cream shadow-sm transition hover:opacity-90"
                    onclick={() => (detailItem = flashCurrent)}
                >
                    Lihat detail
                </button>
            </div>
        </div>
    {/if}

    {#if !loading && !error && filtered.length && view !== "flashcard"}
        <Pagination
            page={currentPage}
            {totalPages}
            onPageChange={(p) => (page = p)}
        />

        <p class="mt-3 text-center text-xs text-navy/50">
            Menampilkan {rangeStart}–{rangeEnd} dari {filtered.length} kanji
        </p>
    {/if}
</div>

<DetailModal item={detailItem} onClose={() => (detailItem = null)} />

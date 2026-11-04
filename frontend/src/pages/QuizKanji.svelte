<script>
    import { onMount } from "svelte";
    import { api } from "../lib/api.js";
    import {
        shuffle,
        uniqueSortedLessons,
        normalizeAnswer,
        neighborItems,
    } from "../lib/quiz.js";

    let items = $state([]);
    let loading = $state(true);
    let error = $state("");

    let phase = $state("setup");
    let answerType = $state("essay");
    let lesson = $state("all");
    let batchSize = $state(10);

    let questions = $state([]);
    let index = $state(0);
    let correctCount = $state(0);
    let results = $state([]);

    let inputKana = $state("");
    let inputArti = $state("");
    let selectedKana = $state(null);
    let selectedArti = $state(null);
    let answered = $state(false);
    let kanaOk = $state(false);
    let artiOk = $state(false);
    let kanaInput = $state(null);

    const lessons = $derived(uniqueSortedLessons(items));
    const current = $derived(questions[index] ?? null);
    const total = $derived(questions.length);
    const accuracy = $derived(
        total ? Math.round((correctCount / total) * 100) : 0,
    );

    const answerOptions = [
        { id: "choice", label: "Pilihan Ganda" },
        { id: "essay", label: "Essay" },
    ];

    const batchOptions = [
        { value: 5, label: "5 soal" },
        { value: 10, label: "10 soal" },
        { value: 20, label: "20 soal" },
        { value: 0, label: "Semua" },
    ];

    onMount(load);

    $effect(() => {
        if (
            phase === "playing" &&
            answerType === "essay" &&
            current &&
            !answered &&
            kanaInput &&
            isDesktop()
        ) {
            kanaInput.focus();
        }
    });

    function isDesktop() {
        if (typeof window === "undefined" || !window.matchMedia) return false;
        return window.matchMedia("(min-width: 1024px) and (hover: hover)")
            .matches;
    }

    async function load() {
        loading = true;
        error = "";
        try {
            items = (await api.listKanji()) || [];
        } catch (e) {
            error = e.message;
        } finally {
            loading = false;
        }
    }

    function lessonPool() {
        if (lesson === "all") return items;
        return items.filter((k) => String(k.lesson) === lesson);
    }

    function makeQuestion(item) {
        if (answerType !== "choice") return { ...item };
        const kanaOptions = shuffle([item, ...neighborItems(item, items, 3)]);
        const artiOptions = shuffle([item, ...neighborItems(item, items, 3)]);
        return { ...item, kanaOptions, artiOptions };
    }

    function start() {
        const pool = shuffle(lessonPool());
        const size =
            batchSize === 0 ? pool.length : Math.min(batchSize, pool.length);
        questions = pool.slice(0, size).map(makeQuestion);
        index = 0;
        correctCount = 0;
        results = [];
        resetTransient();
        phase = questions.length ? "playing" : "setup";
    }

    function resetTransient() {
        inputKana = "";
        inputArti = "";
        selectedKana = null;
        selectedArti = null;
        answered = false;
        kanaOk = false;
        artiOk = false;
    }

    function submit(event) {
        event.preventDefault();
        if (answered) return;

        let userKana = inputKana;
        let userArti = inputArti;

        if (answerType === "choice") {
            if (selectedKana == null || selectedArti == null) return;
            kanaOk = selectedKana === current.id;
            artiOk = selectedArti === current.id;
            userKana =
                current.kanaOptions.find((o) => o.id === selectedKana)
                    ?.furigana ?? "";
            userArti = (
                current.artiOptions.find((o) => o.id === selectedArti)
                    ?.terjemahan || []
            ).join(", ");
        } else {
            if (!inputKana.trim() || !inputArti.trim()) return;
            kanaOk =
                normalizeAnswer(inputKana) ===
                normalizeAnswer(current.furigana);
            artiOk = (current.terjemahan || []).some(
                (t) => normalizeAnswer(t) === normalizeAnswer(inputArti),
            );
        }

        const ok = kanaOk && artiOk;
        answered = true;
        if (ok) correctCount += 1;
        results = [
            ...results,
            { kanji: current, userKana, userArti, correct: ok },
        ];
    }

    function next() {
        if (index + 1 >= questions.length) {
            phase = "done";
            return;
        }
        index += 1;
        resetTransient();
    }

    function choiceState(optionId, selected) {
        if (!answered) return selected === optionId ? "selected" : "idle";
        if (optionId === current.id) return "correct";
        if (optionId === selected) return "wrong";
        return "muted";
    }
</script>

<div class="mx-auto max-w-2xl">
    {#if loading}
        <div class="card flex items-center justify-center p-16 text-navy/60">
            <span class="animate-pulse">Memuat data kanji...</span>
        </div>
    {:else if error}
        <div class="card border-red-300 p-8 text-center text-red-700">
            <p class="font-semibold">Gagal memuat data</p>
            <p class="mt-1 text-sm">{error}</p>
            <button class="btn-primary mt-4" onclick={load}>Coba lagi</button>
        </div>
    {:else if !items.length}
        <div class="card p-10 text-center">
            <p class="text-lg font-semibold text-ink">Belum ada kanji</p>
            <p class="mt-1 text-sm text-navy/60">
                Tambahkan kanji melalui halaman Kelola Kanji terlebih dahulu.
            </p>
        </div>
    {:else if phase === "setup"}
        <div class="mb-5 text-center">
            <h1 class="text-xl font-bold text-ink sm:text-2xl">Tebak Kanji</h1>
            <p class="mt-1 text-sm text-navy/60">
                Tulis cara baca (furigana) dan arti dari kanji
            </p>
        </div>

        <div class="card space-y-5 p-4 sm:p-6">
            <div class="grid gap-4 sm:grid-cols-2">
                <div>
                    <label class="label" for="lesson">Lesson</label>
                    <select id="lesson" class="input" bind:value={lesson}>
                        <option value="all">Semua lesson</option>
                        {#each lessons as l}
                            <option value={l}>Lesson {l}</option>
                        {/each}
                    </select>
                </div>
                <div>
                    <label class="label" for="batch"
                        >Jumlah soal per batch</label
                    >
                    <select id="batch" class="input" bind:value={batchSize}>
                        {#each batchOptions as b}
                            <option value={b.value}>{b.label}</option>
                        {/each}
                    </select>
                </div>
            </div>

            <div>
                <p class="label">Tipe jawaban</p>
                <div class="flex flex-wrap gap-2">
                    {#each answerOptions as a}
                        <button
                            class="rounded-full px-3.5 py-1.5 text-xs font-semibold transition sm:text-sm {answerType ===
                            a.id
                                ? 'bg-sage text-cream shadow-sm'
                                : 'border border-steel/30 bg-white text-navy hover:bg-cream/60'}"
                            onclick={() => (answerType = a.id)}
                        >
                            {a.label}
                        </button>
                    {/each}
                </div>
            </div>

            <div class="rounded-xl bg-sage/10 px-4 py-3 text-sm text-navy/70">
                {#if answerType === "choice"}
                    Pilih <span class="font-semibold text-ink">cara baca</span> dan
                    <span class="font-semibold text-ink">arti</span> yang benar dari
                    pilihan yang tersedia.
                {:else}
                    Jawab dengan menuliskan <span class="font-semibold text-ink"
                        >cara baca</span
                    >
                    dan salah satu <span class="font-semibold text-ink">arti</span> yang
                    benar.
                {/if}
            </div>

            <div
                class="flex flex-wrap items-center justify-between gap-3 border-t border-navy/10 pt-4"
            >
                <p class="text-sm text-navy/60">
                    {lessonPool().length} kanji tersedia
                    {#if lesson !== "all"}<span class="text-navy/40"
                            >(Lesson {lesson})</span
                        >{/if}
                </p>
                <button
                    class="btn-primary"
                    onclick={start}
                    disabled={!lessonPool().length}>Mulai Quiz</button
                >
            </div>
        </div>
    {:else if phase === "playing" && current}
        <div class="mb-4">
            <div
                class="mb-2 flex items-center justify-between text-sm font-semibold text-navy/70"
            >
                <span>Soal {index + 1} / {total}</span>
                <span>Skor: <span class="text-sage">{correctCount}</span></span>
            </div>
            <div class="h-2 w-full overflow-hidden rounded-full bg-navy/10">
                <div
                    class="h-full rounded-full bg-steel transition-all"
                    style="width: {((index + (answered ? 1 : 0)) / total) *
                        100}%"
                ></div>
            </div>
        </div>

        <div class="card">
            <div class="sticky top-0 z-20 rounded-t-2xl bg-navy px-5 py-8 text-center">
                <p
                    class="text-[11px] font-semibold uppercase tracking-widest text-cream/50"
                >
                    Kanji
                </p>
                <p
                    class="mt-3 font-jp text-6xl font-bold text-cream sm:text-7xl"
                >
                    {current.kanji}
                </p>
                <p class="mt-3 text-xs text-cream/50">
                    {current.buku ? `${current.buku} · ` : ''}Lesson {current.lesson}
                </p>
            </div>

            {#if answerType === "choice"}
                <form class="space-y-4 p-4 sm:p-5" onsubmit={submit}>
                    <div>
                        <p class="label">Cara baca (furigana)</p>
                        <div class="grid gap-2.5 sm:grid-cols-2">
                            {#each current.kanaOptions as option, i}
                                {@const state = choiceState(
                                    option.id,
                                    selectedKana,
                                )}
                                <button
                                    type="button"
                                    class="flex items-center gap-3 rounded-xl border px-4 py-3 text-left transition {state ===
                                    'idle'
                                        ? 'border-steel/25 bg-white hover:border-steel hover:bg-cream/60'
                                        : state === 'selected'
                                          ? 'border-steel bg-cream/70'
                                          : state === 'correct'
                                            ? 'border-sage bg-sage/15'
                                            : state === 'wrong'
                                              ? 'border-red-400 bg-red-500/10'
                                              : 'border-transparent bg-white/50 opacity-60'}"
                                    onclick={() =>
                                        !answered && (selectedKana = option.id)}
                                    disabled={answered}
                                >
                                    <span
                                        class="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg text-xs font-bold {state ===
                                        'correct'
                                            ? 'bg-sage text-cream'
                                            : state === 'wrong'
                                              ? 'bg-red-500 text-white'
                                              : 'bg-navy/10 text-navy'}"
                                    >
                                        {"ABCD"[i]}
                                    </span>
                                    <span
                                        class="font-jp text-lg font-semibold text-ink"
                                        >{option.furigana}</span
                                    >
                                    {#if answered && option.id === current.id}
                                        <span class="ml-auto text-sage">✓</span>
                                    {:else if answered &&
                                        option.id === selectedKana}
                                        <span class="ml-auto text-red-500"
                                            >✕</span
                                        >
                                    {/if}
                                </button>
                            {/each}
                        </div>
                    </div>

                    <div>
                        <p class="label">Arti</p>
                        <div class="grid gap-2.5 sm:grid-cols-2">
                            {#each current.artiOptions as option, i}
                                {@const state = choiceState(
                                    option.id,
                                    selectedArti,
                                )}
                                <button
                                    type="button"
                                    class="flex items-center gap-3 rounded-xl border px-4 py-3 text-left transition {state ===
                                    'idle'
                                        ? 'border-steel/25 bg-white hover:border-steel hover:bg-cream/60'
                                        : state === 'selected'
                                          ? 'border-steel bg-cream/70'
                                          : state === 'correct'
                                            ? 'border-sage bg-sage/15'
                                            : state === 'wrong'
                                              ? 'border-red-400 bg-red-500/10'
                                              : 'border-transparent bg-white/50 opacity-60'}"
                                    onclick={() =>
                                        !answered && (selectedArti = option.id)}
                                    disabled={answered}
                                >
                                    <span
                                        class="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg text-xs font-bold {state ===
                                        'correct'
                                            ? 'bg-sage text-cream'
                                            : state === 'wrong'
                                              ? 'bg-red-500 text-white'
                                              : 'bg-navy/10 text-navy'}"
                                    >
                                        {"ABCD"[i]}
                                    </span>
                                    <span class="font-semibold text-ink"
                                        >{(option.terjemahan || []).join(
                                            ", ",
                                        )}</span
                                    >
                                    {#if answered && option.id === current.id}
                                        <span class="ml-auto text-sage">✓</span>
                                    {:else if answered &&
                                        option.id === selectedArti}
                                        <span class="ml-auto text-red-500"
                                            >✕</span
                                        >
                                    {/if}
                                </button>
                            {/each}
                        </div>
                    </div>

                    {#if !answered}
                        <button
                            type="submit"
                            class="btn-primary w-full"
                            disabled={selectedKana == null ||
                                selectedArti == null}
                            >Periksa Jawaban</button
                        >
                    {/if}
                </form>
            {:else}
                <form class="space-y-4 p-4 sm:p-5" onsubmit={submit}>
                    <div>
                        <label class="label" for="kana"
                            >Cara baca (furigana)</label
                        >
                        <input
                            id="kana"
                            bind:this={kanaInput}
                            class="input font-jp {answered
                                ? kanaOk
                                    ? 'border-sage'
                                    : 'border-red-400'
                                : ''}"
                            bind:value={inputKana}
                            disabled={answered}
                            autocomplete="off"
                            placeholder="contoh: わたし"
                        />
                        {#if answered && !kanaOk}
                            <p class="mt-1 text-xs text-red-600">
                                Jawaban yang tepat: <span
                                    class="font-jp font-semibold"
                                    >{current.furigana}</span
                                >
                            </p>
                        {/if}
                    </div>
                    <div>
                        <label class="label" for="arti">Arti</label>
                        <input
                            id="arti"
                            class="input {answered
                                ? artiOk
                                    ? 'border-sage'
                                    : 'border-red-400'
                                : ''}"
                            bind:value={inputArti}
                            disabled={answered}
                            autocomplete="off"
                            placeholder="contoh: saya"
                        />
                        {#if answered && !artiOk}
                            <p class="mt-1 text-xs text-red-600">
                                Jawaban yang tepat: <span class="font-semibold"
                                    >{(current.terjemahan || []).join(
                                        ", ",
                                    )}</span
                                >
                            </p>
                        {/if}
                    </div>

                    {#if !answered}
                        <button
                            type="submit"
                            class="btn-primary w-full"
                            disabled={!inputKana.trim() || !inputArti.trim()}
                            >Periksa Jawaban</button
                        >
                    {/if}
                </form>
            {/if}

            {#if answered}
                {@const ok = kanaOk && artiOk}
                <div
                    class="flex flex-col gap-3 border-t border-navy/10 px-4 py-4 sm:flex-row sm:items-center sm:justify-between sm:px-5"
                >
                    <p
                        class="text-sm font-semibold {ok
                            ? 'text-sage'
                            : 'text-red-600'}"
                    >
                        {ok ? "Benar!" : "Kurang tepat"}
                        <span class="ml-1 font-normal text-navy/70">
                            Jawaban: <span
                                class="font-jp font-semibold text-ink"
                                >{current.furigana}</span
                            >
                            <span class="text-navy/50">—</span>
                            <span class="font-semibold text-ink"
                                >{(current.terjemahan || []).join(", ")}</span
                            >
                        </span>
                    </p>
                    <button class="btn-primary shrink-0" onclick={next}>
                        {index + 1 >= total
                            ? "Lihat Hasil"
                            : "Soal berikutnya →"}
                    </button>
                </div>
            {/if}
        </div>
    {:else if phase === "done"}
        <div class="card p-6 text-center sm:p-8">
            <p
                class="text-sm font-semibold uppercase tracking-wide text-navy/50"
            >
                Batch selesai
            </p>
            <p class="mt-3 text-5xl font-bold text-steel">{accuracy}%</p>
            <p class="mt-2 text-lg font-semibold text-ink">
                Skor {correctCount} / {total}
            </p>
            <p class="mt-1 text-sm text-navy/60">
                {accuracy >= 80
                    ? "Luar biasa! Terus pertahankan."
                    : accuracy >= 50
                      ? "Bagus, tingkatkan lagi ya."
                      : "Jangan menyerah, coba lagi!"}
            </p>

            <div class="mt-6 space-y-1.5 text-left">
                {#each results as r, i}
                    <div
                        class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm {r.correct
                            ? 'bg-sage/10'
                            : 'bg-red-500/10'}"
                    >
                        <span
                            class="font-bold {r.correct
                                ? 'text-sage'
                                : 'text-red-500'}">{r.correct ? "✓" : "✕"}</span
                        >
                        <span class="font-jp text-lg font-semibold text-ink"
                            >{r.kanji.kanji}</span
                        >
                        <span class="font-jp text-navy/60"
                            >{r.kanji.furigana}</span
                        >
                        <span class="ml-auto text-navy/70"
                            >{(r.kanji.terjemahan || []).join(", ")}</span
                        >
                    </div>
                {/each}
            </div>

            <div class="mt-6 flex flex-col gap-2 sm:flex-row sm:justify-center">
                <button class="btn-ghost" onclick={() => (phase = "setup")}
                    >Kembali</button
                >
                <button class="btn-primary" onclick={start}>Ulangi Batch</button
                >
            </div>
        </div>
    {/if}
</div>

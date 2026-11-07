<script>
    import { onMount } from "svelte";
    import { api } from "../lib/api.js";
    import {
        plainText,
        sortLessons,
        lessonLabel,
        targetAnswers,
    } from "../lib/bunpou.js";
    import { shuffle, normalizeAnswer } from "../lib/quiz.js";
    import BunpouSentence from "../components/BunpouSentence.svelte";

    let items = $state([]);
    let loading = $state(true);
    let error = $state("");

    let phase = $state("setup");
    let mode = $state("pertanyaan");
    let lesson = $state("all");
    let batchSize = $state(10);
    let hideKanji = $state(false);

    let qHide = $state(false);

    let questions = $state([]);
    let index = $state(0);
    let correctCount = $state(0);
    let results = $state([]);

    let answerInput = $state("");
    let answered = $state(false);
    let answerOk = $state(false);

    const lessons = $derived(sortLessons(items));
    const current = $derived(questions[index] ?? null);
    const total = $derived(questions.length);
    const accuracy = $derived(
        total ? Math.round((correctCount / total) * 100) : 0,
    );
    const availablePairs = $derived(buildPairs().length);

    const batchOptions = [
        { value: 5, label: "5 soal" },
        { value: 10, label: "10 soal" },
        { value: 20, label: "20 soal" },
        { value: 0, label: "Semua" },
    ];

    const modeOptions = [
        { id: "pertanyaan", label: "Pertanyaan" },
        { id: "jawaban", label: "Jawaban" },
        { id: "acak", label: "Acak" },
    ];

    onMount(load);

    async function load() {
        loading = true;
        error = "";
        try {
            items = (await api.listBunpou()) || [];
        } catch (e) {
            error = e.message;
        } finally {
            loading = false;
        }
    }

    function lessonPool() {
        if (lesson === "all") return items;
        return items.filter((b) => String(b.lesson) === lesson);
    }

    function buildPairs() {
        const out = [];
        for (const b of lessonPool()) {
            for (const pola of b.pola || []) {
                for (const contoh of pola.contoh || []) {
                    if (
                        !contoh?.pertanyaan?.kalimat ||
                        !contoh?.jawaban?.kalimat
                    ) {
                        continue;
                    }
                    out.push({
                        bunpou: b,
                        pola,
                        pertanyaan: contoh.pertanyaan,
                        jawaban: contoh.jawaban,
                    });
                }
            }
        }
        return out;
    }

    function buildQuestion(pair) {
        let side = mode;
        if (mode === "acak")
            side = Math.random() < 0.5 ? "pertanyaan" : "jawaban";
        return { pair, side };
    }

    function start() {
        const pool = shuffle(buildPairs());
        const size =
            batchSize === 0 ? pool.length : Math.min(batchSize, pool.length);
        questions = pool.slice(0, size).map(buildQuestion);
        index = 0;
        correctCount = 0;
        results = [];
        resetTransient();
        phase = questions.length ? "playing" : "setup";
    }

    function resetTransient() {
        answerInput = "";
        answered = false;
        answerOk = false;
        qHide = hideKanji;
    }

    function promptOf(q) {
        return q.side === "pertanyaan" ? q.pair.pertanyaan : q.pair.jawaban;
    }

    function targetOf(q) {
        return q.side === "pertanyaan" ? q.pair.jawaban : q.pair.pertanyaan;
    }

    function answerLabel(kalimat) {
        return plainText(kalimat, "kanji");
    }

    function submit(event) {
        event.preventDefault();
        if (answered || !answerInput.trim()) return;
        const target = targetOf(current);
        const norm = normalizeAnswer(answerInput);
        const targets = targetAnswers(target.kalimat);
        answerOk = targets.some((t) => normalizeAnswer(t) === norm);
        answered = true;
        if (answerOk) correctCount += 1;
        results = [
            ...results,
            { q: current, input: answerInput, correct: answerOk },
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
</script>

<div class="mx-auto max-w-2xl">
    {#if loading}
        <div class="card flex items-center justify-center p-16 text-navy/60">
            <span class="animate-pulse">Memuat data bunpou...</span>
        </div>
    {:else if error}
        <div class="card border-red-300 p-8 text-center text-red-700">
            <p class="font-semibold">Gagal memuat data</p>
            <p class="mt-1 text-sm">{error}</p>
            <button class="btn-primary mt-4" onclick={load}>Coba lagi</button>
        </div>
    {:else if !items.length}
        <div class="card p-10 text-center">
            <p class="text-lg font-semibold text-ink">Belum ada bunpou</p>
            <p class="mt-1 text-sm text-navy/60">
                Tambahkan bunpou melalui halaman Kelola Bunpou terlebih dahulu.
            </p>
        </div>
    {:else if phase === "setup"}
        <div class="mb-5 text-center">
            <h1 class="text-xl font-bold text-ink sm:text-2xl">
                Tebak Bunpou
            </h1>
            <p class="mt-1 text-sm text-navy/60">
                Tebak pertanyaan dari sebuah jawaban atau sebaliknya
            </p>
        </div>

        <div class="card space-y-5 p-4 sm:p-6">
            <div>
                <span class="label">Soal yang ditampilkan</span>
                <div class="flex gap-1 rounded-xl bg-navy/5 p-1">
                    {#each modeOptions as m}
                        <button
                            class="flex-1 rounded-lg px-3 py-2 text-xs font-semibold transition sm:text-sm {mode ===
                            m.id
                                ? 'bg-steel text-cream shadow-sm'
                                : 'text-navy/70 hover:bg-white'}"
                            onclick={() => (mode = m.id)}
                        >
                            {m.label}
                        </button>
                    {/each}
                </div>
                <p class="mt-1.5 text-xs text-navy/50">
                    {mode === "pertanyaan"
                        ? "Ditampilkan pertanyaan, jawab dengan jawabannya."
                        : mode === "jawaban"
                          ? "Ditampilkan jawaban, jawab dengan pertanyaannya."
                          : "Arah soal diacak setiap pertanyaan."}
                </p>
            </div>

            <div class="grid gap-4 sm:grid-cols-2">
                <div>
                    <label class="label" for="lesson">Lesson</label>
                    <select id="lesson" class="input" bind:value={lesson}>
                        <option value="all">Semua lesson</option>
                        {#each lessons as l}
                            <option value={l}>{lessonLabel(l)}</option>
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

            <label
                class="flex cursor-pointer items-center justify-between rounded-xl bg-sage/10 px-4 py-3"
            >
                <span class="text-sm font-semibold text-ink"
                    >Sembunyikan kanji</span
                >
                <input
                    type="checkbox"
                    class="h-5 w-5 accent-steel"
                    bind:checked={hideKanji}
                />
            </label>

            <div
                class="flex flex-wrap items-center justify-between gap-3 border-t border-navy/10 pt-4"
            >
                <p class="text-sm text-navy/60">
                    {availablePairs} soal tersedia
                    {#if lesson !== "all"}<span class="text-navy/40"
                            >({lessonLabel(lesson)})</span
                        >{/if}
                </p>
                <button
                    class="btn-primary"
                    onclick={start}
                    disabled={!availablePairs}>Mulai Quiz</button
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

        <div class="card overflow-hidden">
            <div class="bg-navy px-5 py-8 text-center">
                <p
                    class="text-[11px] font-semibold uppercase tracking-widest text-cream/50"
                >
                    {current.side === "pertanyaan" ? "Pertanyaan" : "Jawaban"}
                </p>
                <BunpouSentence
                    kalimat={promptOf(current).kalimat}
                    showTerjemahan={false}
                    hideKanji={qHide}
                    className="mt-3 block text-2xl font-bold leading-relaxed text-cream sm:text-3xl"
                />
                <label
                    class="mt-4 inline-flex cursor-pointer items-center gap-2 text-xs text-cream/60"
                >
                    <input
                        type="checkbox"
                        class="h-4 w-4 accent-sage"
                        bind:checked={qHide}
                    />
                    Sembunyikan kanji
                </label>
            </div>

            <form class="space-y-4 p-4 sm:p-5" onsubmit={submit}>
                <div>
                    <label class="label" for="answer">
                        {current.side === "pertanyaan"
                            ? "Jawaban (kalimat Jepang)"
                            : "Pertanyaan (kalimat Jepang)"}
                    </label>
                    <input
                        id="answer"
                        class="input font-jp {answered
                            ? answerOk
                                ? 'border-sage'
                                : 'border-red-400'
                            : ''}"
                        bind:value={answerInput}
                        disabled={answered}
                        autocomplete="off"
                        placeholder="tulis kalimat bahasa Jepang"
                    />
                    {#if answered && !answerOk}
                        <p class="mt-2 text-xs text-red-600">
                            Jawaban yang tepat:
                        </p>
                        <BunpouSentence
                            kalimat={targetOf(current).kalimat}
                            terjemahan={targetOf(current).terjemahan}
                            hideKanji={qHide}
                            className="mt-1 text-base font-semibold text-ink"
                        />
                    {/if}
                </div>

                {#if !answered}
                    <button
                        type="submit"
                        class="btn-primary w-full"
                        disabled={!answerInput.trim()}>Periksa Jawaban</button
                    >
                {/if}
            </form>

            {#if answered}
                <div
                    class="flex flex-col gap-3 border-t border-navy/10 px-4 py-4 sm:flex-row sm:items-center sm:justify-between sm:px-5"
                >
                    <div class="min-w-0">
                        <p
                            class="text-sm font-semibold {answerOk
                                ? 'text-sage'
                                : 'text-red-600'}"
                        >
                            {answerOk ? "Benar!" : "Kurang tepat"}
                        </p>
                        <BunpouSentence
                            kalimat={targetOf(current).kalimat}
                            terjemahan={targetOf(current).terjemahan}
                            hideKanji={qHide}
                            className="mt-1 block text-base font-semibold text-ink"
                        />
                    </div>
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
                        class="rounded-lg px-3 py-2 text-sm {r.correct
                            ? 'bg-sage/10'
                            : 'bg-red-500/10'}"
                    >
                        <div class="flex items-center gap-2">
                            <span
                                class="font-bold {r.correct
                                    ? 'text-sage'
                                    : 'text-red-500'}"
                                >{r.correct ? "✓" : "✗"}</span
                            >
                            <span class="font-jp font-semibold text-ink"
                                >{answerLabel(promptOf(r.q).kalimat)}</span
                            >
                            <span class="text-navy/40">→</span>
                            <span class="font-jp text-sage"
                                >{answerLabel(targetOf(r.q).kalimat)}</span
                            >
                        </div>
                        {#if !r.correct}
                            <p class="mt-1 pl-6 text-xs text-navy/60">
                                Jawabanmu: <span class="font-jp"
                                    >{r.input}</span
                                >
                            </p>
                        {/if}
                    </div>
                {/each}
            </div>

            <div class="mt-6 flex flex-col gap-2 sm:flex-row sm:justify-center">
                <button class="btn-ghost" onclick={() => (phase = "setup")}
                    >Kembali</button
                >
                <button class="btn-primary" onclick={start}>Ulangi Batch</button>
            </div>
        </div>
    {/if}
</div>

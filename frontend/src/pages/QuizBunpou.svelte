<script>
    import { onMount } from "svelte";
    import { api } from "../lib/api.js";
    import { plainText, lessonLabel } from "../lib/bunpou.js";
    import { shuffle, normalizeAnswer } from "../lib/quiz.js";
    import BunpouSentence from "../components/BunpouSentence.svelte";

    const kategoriOptions = [
        { value: "jtest", label: "J-Test" },
        { value: "pg", label: "Pilihan ganda" },
        { value: "bintang", label: "Bintang" },
        { value: "dokkai", label: "Dokkai" },
        { value: "acak", label: "Acak (semua jenis)" },
    ];

    const kategoriLabel = {
        jtest: "J-Test",
        pg: "Pilihan ganda",
        bintang: "Bintang",
        dokkai: "Dokkai",
    };

    let lessons = $state([]);
    let lesson = $state("");
    let pool = $state([]);
    let loading = $state(true);
    let loadingQuestions = $state(false);
    let error = $state("");

    let phase = $state("setup");
    let batchSize = $state(10);
    let typeFilter = $state("acak");

    let questions = $state([]);
    let index = $state(0);
    let correctCount = $state(0);
    let results = $state([]);

    let answerInput = $state("");
    let selectedIndex = $state(-1);
    let answered = $state(false);
    let answerOk = $state(false);

    const current = $derived(questions[index] ?? null);
    const total = $derived(questions.length);
    const accuracy = $derived(
        total ? Math.round((correctCount / total) * 100) : 0,
    );

    const available = $derived(
        typeFilter === "acak"
            ? pool
            : pool.filter((soal) => soal.kategori === typeFilter),
    );

    const batchOptions = [
        { value: 5, label: "5 soal" },
        { value: 10, label: "10 soal" },
        { value: 20, label: "20 soal" },
        { value: 0, label: "Semua" },
    ];

    onMount(loadLessons);

    async function loadLessons() {
        loading = true;
        error = "";
        try {
            lessons = ((await api.listLatihanBunpouLessons()) || []).map(
                String,
            );
            if (lessons.length) {
                lesson = lessons[0];
                await fetchQuestions();
            } else {
                pool = [];
            }
        } catch (e) {
            error = e.message;
        } finally {
            loading = false;
        }
    }

    async function fetchQuestions() {
        if (!lesson) {
            pool = [];
            return;
        }
        loadingQuestions = true;
        error = "";
        try {
            pool = (await api.getLatihanBunpou(lesson)) || [];
        } catch (e) {
            error = e.message;
            pool = [];
        } finally {
            loadingQuestions = false;
        }
    }

    async function selectLesson(event) {
        lesson = event.target.value;
        await fetchQuestions();
    }

    function start() {
        if (!available.length) return;
        const chosen = shuffle([...available]);
        const size =
            batchSize === 0 ? chosen.length : Math.min(batchSize, chosen.length);
        questions = chosen.slice(0, size);
        index = 0;
        correctCount = 0;
        results = [];
        resetTransient();
        phase = "playing";
    }

    function resetTransient() {
        answerInput = "";
        selectedIndex = -1;
        answered = false;
        answerOk = false;
    }

    function isObjective(soal) {
        return soal?.kategori !== "jtest";
    }

    function answerIndex(value) {
        const text = String(value ?? "").trim().toLowerCase();
        if (text.length !== 1) return -1;
        const code = text.charCodeAt(0) - 97;
        return code >= 0 && code <= 25 ? code : -1;
    }

    function correctAnswerText(soal) {
        if (!soal) return "";
        if (!isObjective(soal)) {
            return soal.jawaban_benar || soal.jawaban || "";
        }
        const idx = answerIndex(soal.jawaban_benar);
        const option = soal.pilihan_jawaban?.[idx];
        return option ?? soal.jawaban_benar ?? "";
    }

    function isCorrectOption(optionIndex) {
        if (!current) return false;
        return answerIndex(current.jawaban_benar) === optionIndex;
    }

    function chooseOption(optionIndex) {
        if (answered || !current) return;
        selectedIndex = optionIndex;
        answerOk = isCorrectOption(optionIndex);
        finishAnswer(String.fromCharCode(65 + optionIndex));
    }

    function submitEssay(event) {
        event.preventDefault();
        if (answered || !answerInput.trim() || !current) return;
        const expected = current.jawaban_benar || current.jawaban || "";
        const variants = [
            expected,
            plainText(expected, "kanji"),
            plainText(expected, "kana"),
        ].filter(Boolean);
        const norm = normalizeAnswer(answerInput);
        answerOk = variants.some((item) => normalizeAnswer(item) === norm);
        finishAnswer(answerInput);
    }

    function finishAnswer(input) {
        answered = true;
        if (answerOk) correctCount += 1;
        results = [...results, { q: current, input, correct: answerOk }];
    }

    function next() {
        if (index + 1 >= questions.length) {
            phase = "done";
            return;
        }
        index += 1;
        resetTransient();
    }

    function answerLabel(value) {
        return plainText(value, "kanji");
    }

    function optionLetter(optionIndex) {
        return String.fromCharCode(65 + optionIndex);
    }

    function preserveSpaces(soal) {
        return soal?.kategori === "bintang" || soal?.kategori === "jtest";
    }

    function typeHint(soal) {
        if (!soal) return "";
        switch (soal.kategori) {
            case "jtest":
                return soal.instruksi || "Susunlah jawaban yang benar.";
            case "bintang":
                return "Pilih kata yang tepat untuk posisi ★.";
            case "dokkai":
                return "Bacalah dokkai lalu jawab pertanyaan.";
            default:
                return "Pilih jawaban yang benar.";
        }
    }
</script>

<div class="mx-auto max-w-2xl">
    {#if loading}
        <div class="card flex items-center justify-center p-16 text-navy/60">
            <span class="animate-pulse">Memuat data latihan...</span>
        </div>
    {:else if error}
        <div class="card border-red-300 p-8 text-center text-red-700">
            <p class="font-semibold">Gagal memuat data</p>
            <p class="mt-1 text-sm">{error}</p>
            <button class="btn-primary mt-4" onclick={loadLessons}
                >Coba lagi</button
            >
        </div>
    {:else if !lessons.length}
        <div class="card p-10 text-center">
            <p class="text-lg font-semibold text-ink">Belum ada latihan bunpou</p>
            <p class="mt-1 text-sm text-navy/60">
                Tambahkan file soal di
                <span class="font-mono">data/latihan/bunpou/lesson1</span>.
            </p>
        </div>
    {:else if phase === "setup"}
        <div class="mb-5 text-center">
            <h1 class="text-xl font-bold text-ink sm:text-2xl">Latihan Bunpou</h1>
            <p class="mt-1 text-sm text-navy/60">
                Pilih lesson dan jenis soal yang ingin dikerjakan.
            </p>
        </div>

        <div class="card space-y-5 p-4 sm:p-6">
            <div class="grid gap-4 sm:grid-cols-2">
                <div>
                    <label class="label" for="lesson">Lesson</label>
                    <select
                        id="lesson"
                        class="input"
                        bind:value={lesson}
                        onchange={selectLesson}
                    >
                        {#each lessons as l}
                            <option value={l}>{lessonLabel(l)}</option>
                        {/each}
                    </select>
                </div>
                <div>
                    <label class="label" for="jenis">Jenis soal</label>
                    <select id="jenis" class="input" bind:value={typeFilter}>
                        {#each kategoriOptions as option}
                            <option value={option.value}>{option.label}</option>
                        {/each}
                    </select>
                </div>
                <div class="sm:col-span-2">
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

            <div
                class="flex flex-wrap items-center justify-between gap-3 border-t border-navy/10 pt-4"
            >
                <p class="text-sm text-navy/60">
                    {#if loadingQuestions}
                        Memuat soal...
                    {:else}
                        {available.length} soal tersedia
                    {/if}
                </p>
                <button
                    class="btn-primary"
                    onclick={start}
                    disabled={!available.length || loadingQuestions}
                    >Mulai Latihan</button
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
            <div class="bg-navy px-5 py-7 text-center">
                <div
                    class="flex flex-wrap items-center justify-center gap-2 text-[11px] font-semibold uppercase tracking-widest text-cream/50"
                >
                    <span>{kategoriLabel[current.kategori] || current.kategori}</span>
                    <span class="text-cream/25">•</span>
                    <span>{lessonLabel(current.lesson)}</span>
                    {#if current.no}
                        <span class="text-cream/25">•</span>
                        <span>No. {current.no}</span>
                    {/if}
                </div>

                <p class="mt-1 text-xs text-cream/60">{typeHint(current)}</p>

                <BunpouSentence
                    kalimat={current.soal}
                    showTerjemahan={false}
                    preserveSpaces={preserveSpaces(current)}
                    className="mt-3 block text-2xl font-bold leading-relaxed text-cream sm:text-3xl"
                />
            </div>

            {#if current.kategori === "dokkai" && current.dokkai}
                <div class="border-b border-navy/10 bg-sand/25 px-4 py-3 sm:px-5">
                    <p
                        class="mb-1.5 text-[10px] font-bold uppercase tracking-widest text-steel/70"
                    >
                        Dokkai
                    </p>
                    <BunpouSentence
                        kalimat={current.dokkai}
                        showTerjemahan={false}
                        className="block text-sm leading-relaxed text-navy/80"
                    />
                </div>
            {/if}

            {#if isObjective(current)}
                <div class="grid gap-2 p-4 sm:grid-cols-2 sm:p-5">
                    {#each current.pilihan_jawaban || [] as option, optionIndex}
                        <button
                            type="button"
                            class="flex w-full items-start gap-2.5 rounded-xl border px-3 py-2.5 text-left transition {answered
                                ? isCorrectOption(optionIndex)
                                    ? 'border-sage bg-sage/15'
                                    : selectedIndex === optionIndex
                                      ? 'border-red-400 bg-red-500/10'
                                      : 'border-steel/20 bg-white opacity-60'
                                : 'border-steel/25 bg-white hover:border-steel/50 hover:bg-cream/60'}"
                            onclick={() => chooseOption(optionIndex)}
                            disabled={answered}
                        >
                            <span
                                class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-md bg-steel/10 text-[10px] font-bold text-navy/60"
                                >{optionLetter(optionIndex)}</span
                            >
                            <BunpouSentence
                                kalimat={option}
                                showTerjemahan={false}
                                preserveSpaces={current.kategori === "bintang"}
                                className="text-base font-semibold text-ink"
                            />
                        </button>
                    {/each}
                </div>
            {:else}
                <form class="space-y-4 p-4 sm:p-5" onsubmit={submitEssay}>
                    <div>
                        <label class="label" for="answer"
                            >Jawaban (tulis dalam bahasa Jepang)</label
                        >
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
                            placeholder="tulis jawaban bahasa Jepang"
                        />
                    </div>

                    {#if !answered}
                        <button
                            type="submit"
                            class="btn-primary w-full"
                            disabled={!answerInput.trim()}
                            >Periksa Jawaban</button
                        >
                    {/if}
                </form>
            {/if}

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
                        {#if !answerOk || isObjective(current)}
                            <p class="mt-1 text-xs text-navy/60">
                                Jawaban yang tepat:
                                <span class="font-jp font-semibold text-ink"
                                    >{answerLabel(
                                        correctAnswerText(current),
                                    )}</span
                                >
                            </p>
                        {/if}
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
                        <div class="flex flex-wrap items-center gap-2">
                            <span
                                class="font-bold {r.correct
                                    ? 'text-sage'
                                    : 'text-red-500'}"
                                >{r.correct ? "✓" : "✗"}</span
                            >
                            <span
                                class="rounded bg-sand/50 px-1.5 py-0.5 text-[10px] font-semibold uppercase tracking-wide text-navy/60"
                                >{kategoriLabel[r.q.kategori] ||
                                    r.q.kategori}</span
                            >
                            <span class="font-jp font-semibold text-ink"
                                >{answerLabel(r.q.soal)}</span
                            >
                            <span class="text-navy/40">→</span>
                            <span class="font-jp text-sage"
                                >{answerLabel(correctAnswerText(r.q))}</span
                            >
                        </div>
                        {#if !r.correct}
                            <p class="mt-1 pl-6 text-xs text-navy/60">
                                Jawabanmu: <span class="font-jp">{r.input}</span
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

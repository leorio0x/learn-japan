<script>
    import { onMount } from "svelte";
    import { api } from "../lib/api.js";
    import {
        shuffle,
        uniqueSortedLessons,
        normalizeAnswer,
        isDateWord,
        neighborItems,
    } from "../lib/quiz.js";
    import Furigana from "../components/Furigana.svelte";

    let words = $state([]);
    let loading = $state(true);
    let error = $state("");

    let phase = $state("setup");
    let direction = $state("random");
    let answerType = $state("choice");
    let selectedLessons = $state([]);
    let rangeInput = $state("");
    let rangeError = $state("");
    let batchSize = $state(10);
    let hideKanji = $state(false);

    let questions = $state([]);
    let index = $state(0);
    let correctCount = $state(0);
    let results = $state([]);

    let selectedId = $state(null);
    let essayAnswer = $state("");
    let answered = $state(false);

    const lessons = $derived(uniqueSortedLessons(words));
    const lessonGroups = $derived.by(() => {
        const bookOf = new Map();
        for (const w of words) {
            const l = String(w.lesson);
            if (!bookOf.has(l)) bookOf.set(l, w.buku || "");
        }
        const groups = new Map();
        for (const l of lessons) {
            const buku = bookOf.get(String(l)) || "";
            if (!groups.has(buku)) groups.set(buku, []);
            groups.get(buku).push(l);
        }
        return [...groups.entries()].map(([buku, list]) => ({
            buku,
            label: bookLabel(buku),
            lessons: list,
        }));
    });
    const lessonSummary = $derived.by(() => {
        if (!selectedLessons.length) return "";
        const sorted = [...selectedLessons].sort((a, b) => Number(a) - Number(b));
        return sorted.length <= 8
            ? sorted.join(", ")
            : `${sorted.slice(0, 8).join(", ")}…`;
    });
    const current = $derived(questions[index] ?? null);
    const autoHideKanji = $derived(isDateWord(current?.word));
    const effectiveHideKanji = $derived(hideKanji || autoHideKanji);
    const total = $derived(questions.length);
    const accuracy = $derived(
        total ? Math.round((correctCount / total) * 100) : 0,
    );

    const directionOptions = [
        { id: "jp-id", label: "日本語 → Indonesia" },
        { id: "id-jp", label: "Indonesia → 日本語" },
        { id: "random", label: "Acak" },
    ];

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

    async function load() {
        loading = true;
        error = "";
        try {
            words = (await api.list()) || [];
        } catch (e) {
            error = e.message;
        } finally {
            loading = false;
        }
    }

    function lessonPool() {
        if (!selectedLessons.length) return words;
        const chosen = new Set(selectedLessons);
        return words.filter((w) => chosen.has(String(w.lesson)));
    }

    function toggleLesson(value) {
        const key = String(value);
        rangeError = "";
        selectedLessons = selectedLessons.includes(key)
            ? selectedLessons.filter((l) => l !== key)
            : [...selectedLessons, key];
    }

    function resetLessons() {
        rangeError = "";
        selectedLessons = [];
    }

    function bookLabel(buku) {
        const value = String(buku || "").toLowerCase();
        if (value.includes("1")) return "Minna no Nihongo 1";
        if (value.includes("2")) return "Minna no Nihongo 2";
        return buku || "Lainnya";
    }

    function parseLessonRange(input) {
        const tokens = input
            .split(",")
            .map((t) => t.trim())
            .filter(Boolean);
        if (!tokens.length) return { error: "Isi rentang lesson dulu." };
        const picked = new Set();
        for (const token of tokens) {
            const range = token.match(/^(\d+)\s*-\s*(\d+)$/);
            if (range) {
                let start = Number(range[1]);
                let end = Number(range[2]);
                if (start > end) [start, end] = [end, start];
                for (let i = start; i <= end; i++) picked.add(String(i));
            } else if (/^\d+$/.test(token)) {
                picked.add(String(Number(token)));
            } else {
                return { error: `Format "${token}" tidak valid.` };
            }
        }
        return { picked: [...picked] };
    }

    function applyRange() {
        const result = parseLessonRange(rangeInput);
        if (result.error) {
            rangeError = result.error;
            return;
        }
        const available = new Set(lessons.map(String));
        const valid = result.picked
            .filter((l) => available.has(l))
            .sort((a, b) => Number(a) - Number(b));
        if (!valid.length) {
            rangeError = "Tidak ada lesson yang cocok.";
            return;
        }
        rangeError = "";
        selectedLessons = valid;
    }

    function makeQuestion(word) {
        const dir =
            direction === "random"
                ? Math.random() < 0.5
                    ? "jp-id"
                    : "id-jp"
                : direction;
        let options = null;
        if (answerType === "choice") {
            const distractors = neighborItems(word, words, 3);
            options = shuffle([word, ...distractors]);
        }
        return { word, dir, options };
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
        selectedId = null;
        essayAnswer = "";
        answered = false;
    }

    function record(ok, userAnswer) {
        answered = true;
        if (ok) correctCount += 1;
        results = [
            ...results,
            { word: current.word, dir: current.dir, userAnswer, correct: ok },
        ];
    }

    function choose(option) {
        if (answered) return;
        selectedId = option.id;
        record(option.id === current.word.id, option.id);
    }

    function submitEssay(event) {
        event.preventDefault();
        if (answered || !essayAnswer.trim()) return;
        record(isEssayCorrect(), essayAnswer);
    }

    function isEssayCorrect() {
        const w = current.word;
        const given = normalizeAnswer(essayAnswer);
        if (current.dir === "jp-id")
            return (w.terjemahan || []).some(
                (t) => normalizeAnswer(t) === given,
            );
        return (
            given === normalizeAnswer(w.kanji) ||
            given === normalizeAnswer(w.furigana)
        );
    }

    function next() {
        if (index + 1 >= questions.length) {
            phase = "done";
            return;
        }
        index += 1;
        resetTransient();
    }

    function correctAnswerText(item) {
        const w = item.word;
        if (item.dir === "jp-id") return (w.terjemahan || []).join(", ");
        if (effectiveHideKanji) return w.furigana;
        return w.kanji
            ? `${w.kanji} (${w.furigana})`
            : w.furigana;
    }

    function answerState(option) {
        if (!answered) return "idle";
        if (option.id === current.word.id) return "correct";
        if (option.id === selectedId) return "wrong";
        return "muted";
    }
</script>

<div class="mx-auto max-w-2xl">
    {#if loading}
        <div class="card flex items-center justify-center p-16 text-navy/60">
            <span class="animate-pulse">Memuat data kotoba...</span>
        </div>
    {:else if error}
        <div class="card border-red-300 p-8 text-center text-red-700">
            <p class="font-semibold">Gagal memuat data</p>
            <p class="mt-1 text-sm">{error}</p>
            <button class="btn-primary mt-4" onclick={load}>Coba lagi</button>
        </div>
    {:else if !words.length}
        <div class="card p-10 text-center">
            <p class="text-lg font-semibold text-ink">Belum ada kotoba</p>
            <p class="mt-1 text-sm text-navy/60">
                Tambahkan kosakata melalui halaman Kelola Kotoba terlebih
                dahulu.
            </p>
        </div>
    {:else if phase === "setup"}
        <div class="mb-5 text-center">
            <h1 class="text-xl font-bold text-ink sm:text-2xl">Tebak Kotoba</h1>
            <p class="mt-1 text-sm text-navy/60">
                Atur mode quiz lalu mulai batch pertanyaan
            </p>
        </div>

        <div class="card space-y-5 p-4 sm:p-6">
            <div>
                <p class="label">Arah terjemahan</p>
                <div class="flex flex-wrap gap-2">
                    {#each directionOptions as d}
                        <button
                            class="rounded-full px-3.5 py-1.5 text-xs font-semibold transition sm:text-sm {direction ===
                            d.id
                                ? 'bg-steel text-cream shadow-sm'
                                : 'border border-steel/30 bg-white text-navy hover:bg-cream/60'}"
                            onclick={() => (direction = d.id)}
                        >
                            {d.label}
                        </button>
                    {/each}
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

            <div>
                <div class="flex items-center justify-between gap-2">
                    <p class="label mb-0">Lesson</p>
                    {#if selectedLessons.length}
                        <button
                            type="button"
                            class="text-xs font-semibold text-steel hover:underline"
                            onclick={resetLessons}>Pilih semua</button
                        >
                    {/if}
                </div>

                <div class="mt-2 max-h-64 space-y-3 overflow-y-auto pr-1">
                    {#each lessonGroups as group}
                        <div>
                            <p
                                class="mb-1.5 text-[11px] font-semibold uppercase tracking-wide text-navy/40"
                            >
                                {group.label}
                            </p>
                            <div class="grid grid-cols-4 gap-2">
                                {#each group.lessons as l}
                                    <button
                                        type="button"
                                        class="rounded-full px-2 py-1.5 text-center text-xs font-semibold leading-tight transition {selectedLessons.includes(
                                            String(l),
                                        )
                                            ? 'bg-steel text-cream shadow-sm'
                                            : 'border border-steel/30 bg-white text-navy hover:bg-cream/60'}"
                                        onclick={() => toggleLesson(l)}
                                    >
                                        Lesson {l}
                                    </button>
                                {/each}
                            </div>
                        </div>
                    {/each}
                </div>

                <div class="mt-3 flex flex-col gap-2 sm:flex-row sm:items-center">
                    <input
                        class="input text-sm"
                        bind:value={rangeInput}
                        placeholder="Rentang cepat: 1-10, 15"
                        onkeydown={(e) => e.key === "Enter" && applyRange()}
                    />
                    <button
                        type="button"
                        class="btn-ghost shrink-0 text-xs"
                        onclick={applyRange}>Terapkan</button
                    >
                </div>
                {#if rangeError}
                    <p class="mt-1.5 text-xs text-red-600">{rangeError}</p>
                {/if}
                <p class="mt-2 text-xs text-navy/50">
                    {selectedLessons.length
                        ? `${selectedLessons.length} lesson dipilih`
                        : "Semua lesson akan digunakan"}
                </p>
            </div>

            <div>
                <label class="label" for="batch">Jumlah soal per batch</label>
                <select id="batch" class="input" bind:value={batchSize}>
                    {#each batchOptions as b}
                        <option value={b.value}>{b.label}</option>
                    {/each}
                </select>
            </div>

            <div
                class="flex flex-wrap items-center justify-between gap-3 border-t border-navy/10 pt-4"
            >
                <p class="text-sm text-navy/60">
                    {lessonPool().length} kotoba tersedia
                    {#if selectedLessons.length}
                        <span class="text-navy/40"
                            >(Lesson {lessonSummary})</span
                        >
                    {/if}
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

        <div class="mb-3 flex flex-col items-end gap-1">
            <button
                type="button"
                class="flex items-center gap-2 rounded-full border border-steel/30 bg-white/70 px-3.5 py-1.5 text-xs font-semibold text-navy transition hover:bg-white disabled:cursor-not-allowed disabled:opacity-80"
                onclick={() => (hideKanji = !hideKanji)}
                aria-pressed={effectiveHideKanji}
                disabled={autoHideKanji}
            >
                <span
                    class="relative inline-flex h-5 w-9 items-center rounded-full transition {effectiveHideKanji
                        ? 'bg-sage'
                        : 'bg-navy/20'}"
                >
                    <span
                        class="absolute h-4 w-4 rounded-full bg-white shadow transition-all {effectiveHideKanji
                            ? 'left-4'
                            : 'left-0.5'}"
                    ></span>
                </span>
                Sembunyikan Kanji
                <span class="font-jp text-steel"
                    >({effectiveHideKanji ? "かな" : "漢字"})</span
                >
            </button>
            {#if autoHideKanji}
                <p class="text-[11px] text-navy/50">
                    Soal mengandung tanggal/angka, kanji otomatis
                    disembunyikan.
                </p>
            {/if}
        </div>

        <div class="card overflow-hidden">
            <div class="bg-navy px-5 py-6 text-center sm:py-8">
                <p
                    class="text-[11px] font-semibold uppercase tracking-widest text-cream/50"
                >
                    {current.dir === "jp-id"
                        ? "Bahasa Jepang"
                        : "Bahasa Indonesia"}
                </p>
                <div
                    class="mt-3 flex min-h-[4.5rem] items-center justify-center"
                >
                    {#if current.dir === "jp-id"}
                        <Furigana
                            kanji={current.word.kanji}
                            kana={current.word.furigana}
                            hideKanji={effectiveHideKanji}
                            className="font-jp text-5xl font-bold text-cream sm:text-6xl"
                        />
                    {:else}
                        <p class="text-3xl font-bold text-cream sm:text-4xl">
                            {(current.word.terjemahan || []).join(", ")}
                        </p>
                    {/if}
                </div>
                <p class="mt-3 text-xs text-cream/50">
                    {current.word.buku} &middot; Lesson {current.word.lesson}
                </p>
            </div>

            {#if answerType === "choice"}
                <div class="grid gap-2.5 p-4 sm:grid-cols-2 sm:p-5">
                    {#each current.options as option, i}
                        {@const state = answerState(option)}
                        <button
                            class="flex items-center gap-3 rounded-xl border px-4 py-3 text-left transition {state ===
                            'idle'
                                ? 'border-steel/25 bg-white hover:border-steel hover:bg-cream/60'
                                : state === 'correct'
                                  ? 'border-sage bg-sage/15'
                                  : state === 'wrong'
                                    ? 'border-red-400 bg-red-500/10'
                                    : 'border-transparent bg-white/50 opacity-60'}"
                            onclick={() => choose(option)}
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
                            {#if current.dir === "jp-id"}
                                <span class="text-base font-semibold text-ink"
                                    >{(option.terjemahan || []).join(", ")}</span
                                >
                            {:else}
                                <Furigana
                                    kanji={option.kanji}
                                    kana={option.furigana}
                                    hideKanji={effectiveHideKanji}
                                    className="font-jp text-lg font-semibold text-ink"
                                />
                            {/if}
                            {#if answered && option.id === current.word.id}
                                <span class="ml-auto text-sage">✓</span>
                            {:else if answered && option.id === selectedId}
                                <span class="ml-auto text-red-500">✕</span>
                            {/if}
                        </button>
                    {/each}
                </div>
            {:else}
                <form class="p-4 sm:p-5" onsubmit={submitEssay}>
                    <label class="label" for="essay">
                        {current.dir === "jp-id"
                            ? "Tulis arti dalam bahasa Indonesia"
                            : "Tulis aksara Jepang (kanji / furigana)"}
                    </label>
                    <input
                        id="essay"
                        class="input {current.dir === 'id-jp'
                            ? 'font-jp'
                            : ''} {answered
                            ? results.at(-1)?.correct
                                ? 'border-sage'
                                : 'border-red-400'
                            : ''}"
                        bind:value={essayAnswer}
                        disabled={answered}
                        autocomplete="off"
                        placeholder={current.dir === "jp-id"
                            ? "contoh: saya"
                            : "contoh: 私 / わたし"}
                    />
                    {#if !answered}
                        <button
                            type="submit"
                            class="btn-primary mt-3 w-full"
                            disabled={!essayAnswer.trim()}
                            >Periksa Jawaban</button
                        >
                    {/if}
                </form>
            {/if}

            {#if answered}
                {@const ok = results.at(-1)?.correct}
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
                                >{correctAnswerText(current)}</span
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
                        <span class="font-jp font-semibold text-ink"
                            >{r.word.kanji || r.word.furigana}</span
                        >
                        <span class="text-navy/60"
                            >({r.word.furigana})</span
                        >
                        <span class="ml-auto text-navy/70"
                            >{(r.word.terjemahan || []).join(", ")}</span
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

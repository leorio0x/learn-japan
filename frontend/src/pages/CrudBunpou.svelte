<script>
    import { onMount } from "svelte";
    import { api } from "../lib/api.js";
    import { lessonLabel } from "../lib/bunpou.js";
    import BunpouSentence from "../components/BunpouSentence.svelte";
    import Pagination from "../components/Pagination.svelte";

    const tipeOptions = ["pilihan ganda", "essay"];

    let items = $state([]);
    let lessons = $state([]);
    let loading = $state(true);
    let error = $state("");
    let search = $state("");
    let lessonFilter = $state("all");

    let showForm = $state(false);
    let editingId = $state(null);
    let editingLesson = $state("");
    let saving = $state(false);
    let formError = $state("");
    let form = $state(emptyForm());

    let deleteTarget = $state(null);
    let deleting = $state(false);

    const filtered = $derived(
        items.filter((item) => {
            if (
                lessonFilter !== "all" &&
                String(item.lesson) !== String(lessonFilter)
            ) {
                return false;
            }
            const q = search.trim().toLowerCase();
            if (!q) return true;
            return [
                item.soal,
                item.jawaban_benar,
                ...(item.pilihan_jawaban || []),
            ]
                .filter(Boolean)
                .some((value) => String(value).toLowerCase().includes(q));
        }),
    );

    const pageSize = 10;
    let page = $state(1);
    const totalPages = $derived(
        Math.max(1, Math.ceil(filtered.length / pageSize)),
    );
    const currentPage = $derived(Math.min(page, totalPages));
    const paged = $derived(
        filtered.slice((currentPage - 1) * pageSize, currentPage * pageSize),
    );
    const rangeStart = $derived((currentPage - 1) * pageSize + 1);
    const rangeEnd = $derived(
        Math.min(currentPage * pageSize, filtered.length),
    );

    const choiceOptions = $derived(
        (form.pilihan_jawaban || [])
            .map((choice) => choice.trim())
            .filter(Boolean),
    );

    onMount(load);

    function emptyForm() {
        return {
            lesson: "1",
            soal: "",
            tipe_soal: "pilihan ganda",
            pilihan_jawaban: ["", "", "", ""],
            jawaban_benar: "",
        };
    }

    async function load() {
        loading = true;
        error = "";
        try {
            lessons = (
                (await api.listLatihanBunpouLessons()) || []
            ).map(String);
            const collections = await Promise.all(
                lessons.map(async (lesson) => {
                    const list = (await api.getLatihanBunpou(lesson)) || [];
                    return list.map((item) => ({ ...item, lesson }));
                }),
            );
            items = collections.flat();
            if (form.lesson === "1" && lessons.length) {
                form.lesson = lessons[0];
            }
        } catch (e) {
            error = e.message;
        } finally {
            loading = false;
        }
    }

    function openCreate() {
        form = emptyForm();
        form.lesson =
            lessonFilter !== "all"
                ? String(lessonFilter)
                : lessons[0] || "1";
        editingId = null;
        editingLesson = "";
        formError = "";
        showForm = true;
    }

    function openEdit(item) {
        form = {
            lesson: String(item.lesson),
            soal: item.soal || "",
            tipe_soal: item.tipe_soal || "pilihan ganda",
            pilihan_jawaban:
                item.pilihan_jawaban?.length > 0
                    ? [...item.pilihan_jawaban]
                    : ["", "", "", ""],
            jawaban_benar: item.jawaban_benar || "",
        };
        editingId = item.id;
        editingLesson = String(item.lesson);
        formError = "";
        showForm = true;
    }

    function closeForm() {
        showForm = false;
        editingId = null;
        editingLesson = "";
        formError = "";
        form = emptyForm();
    }

    function addChoice() {
        form.pilihan_jawaban = [...form.pilihan_jawaban, ""];
    }

    function removeChoice(index) {
        form.pilihan_jawaban = form.pilihan_jawaban.filter(
            (_, i) => i !== index,
        );
    }

    async function submit(event) {
        event.preventDefault();
        formError = "";

        const soal = form.soal.trim();
        const jawaban = form.jawaban_benar.trim();
        const choices = form.pilihan_jawaban
            .map((choice) => choice.trim())
            .filter(Boolean);

        if (!soal) {
            formError = "Soal wajib diisi";
            return;
        }
        if (!jawaban) {
            formError = "Jawaban benar wajib diisi";
            return;
        }
        if (form.tipe_soal === "pilihan ganda") {
            if (choices.length < 2) {
                formError = "Pilihan ganda minimal memiliki 2 pilihan jawaban";
                return;
            }
            if (!choices.includes(jawaban)) {
                formError =
                    "Jawaban benar harus salah satu dari pilihan jawaban";
                return;
            }
        }

        const payload = {
            lesson: String(form.lesson || "").trim(),
            soal,
            tipe_soal: form.tipe_soal,
            pilihan_jawaban:
                form.tipe_soal === "pilihan ganda" ? choices : [],
            jawaban_benar: jawaban,
        };

        saving = true;
        try {
            if (editingId) {
                await api.updateLatihanBunpou(
                    editingLesson,
                    editingId,
                    payload,
                );
            } else {
                await api.createLatihanBunpou(payload);
            }
            await load();
            closeForm();
        } catch (e) {
            formError = e.message;
        } finally {
            saving = false;
        }
    }

    async function confirmDelete() {
        if (!deleteTarget) return;
        deleting = true;
        error = "";
        try {
            await api.removeLatihanBunpou(
                deleteTarget.lesson,
                deleteTarget.id,
            );
            deleteTarget = null;
            await load();
        } catch (e) {
            error = e.message;
        } finally {
            deleting = false;
        }
    }

    function isCorrectChoice(item, choice) {
        return String(item.jawaban_benar) === String(choice);
    }
</script>

<div class="mx-auto max-w-3xl">
    <div
        class="mb-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
    >
        <div>
            <h1 class="text-xl font-bold text-ink sm:text-2xl">Kelola Bunpou</h1>
            <p class="mt-1 text-sm text-navy/60">
                Kelola soal latihan bunpou sesuai format soal_bunpou
            </p>
        </div>
        {#if !showForm}
            <button
                class="btn-primary w-full sm:w-auto"
                onclick={openCreate}>+ Tambah Soal</button
            >
        {/if}
    </div>

    {#if showForm}
        <form class="card mb-5 p-4 sm:p-5" onsubmit={submit}>
            <div class="mb-4 flex items-center justify-between">
                <h2 class="font-bold text-ink">
                    {editingId ? "Ubah Soal" : "Soal Baru"}
                </h2>
                <button
                    type="button"
                    class="text-sm font-medium text-navy/50 hover:text-navy"
                    onclick={closeForm}>Tutup</button
                >
            </div>

            <div class="grid gap-3.5 sm:grid-cols-2">
                <div>
                    <label class="label" for="lesson">Lesson</label>
                    <input
                        id="lesson"
                        class="input"
                        type="number"
                        min="1"
                        bind:value={form.lesson}
                        disabled={!!editingId}
                        placeholder="1"
                    />
                </div>
                <div>
                    <label class="label" for="tipe">Tipe soal</label>
                    <select id="tipe" class="input" bind:value={form.tipe_soal}>
                        {#each tipeOptions as tipe}
                            <option value={tipe}>{tipe}</option>
                        {/each}
                    </select>
                </div>
            </div>

            <div class="mt-3">
                <label class="label" for="soal">Soal</label>
                <textarea
                    id="soal"
                    class="input input-textarea font-jp"
                    rows="3"
                    bind:value={form.soal}
                    placeholder="contoh: <私:わたし:saya>（　）<学生:がくせい:mahasiswa> です。"
                ></textarea>
            </div>

            {#if form.tipe_soal === "pilihan ganda"}
                <div class="mt-3">
                    <div class="mb-2 flex items-center justify-between">
                        <span class="label mb-0">Pilihan jawaban</span>
                        <button
                            type="button"
                            class="text-xs font-semibold text-steel hover:text-navy"
                            onclick={addChoice}>+ Tambah pilihan</button
                        >
                    </div>
                    <div class="space-y-2">
                        {#each form.pilihan_jawaban as _, index}
                            <div class="flex items-center gap-2">
                                <span
                                    class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-steel/10 text-xs font-bold text-navy/60"
                                    >{String.fromCharCode(65 + index)}</span
                                >
                                <input
                                    class="input font-jp"
                                    bind:value={form.pilihan_jawaban[index]}
                                    placeholder="pilihan jawaban"
                                />
                                <button
                                    type="button"
                                    class="shrink-0 rounded-lg border border-red-300 px-2.5 py-1.5 text-xs font-bold text-red-500 transition hover:bg-red-500/10 disabled:opacity-40"
                                    onclick={() => removeChoice(index)}
                                    disabled={form.pilihan_jawaban.length <= 2}
                                    aria-label="Hapus pilihan"
                                    >✕</button
                                >
                            </div>
                        {/each}
                    </div>
                    <div class="mt-3">
                        <label class="label" for="jawaban-pilihan"
                            >Jawaban benar</label
                        >
                        <select
                            id="jawaban-pilihan"
                            class="input"
                            bind:value={form.jawaban_benar}
                        >
                            <option value="">Pilih jawaban benar</option>
                            {#each choiceOptions as choice}
                                <option value={choice}>{choice}</option>
                            {/each}
                        </select>
                    </div>
                </div>
            {:else}
                <div class="mt-3">
                    <label class="label" for="jawaban-essay">Jawaban benar</label>
                    <input
                        id="jawaban-essay"
                        class="input font-jp"
                        bind:value={form.jawaban_benar}
                        placeholder="tulis jawaban essay"
                    />
                </div>
            {/if}

            {#if formError}
                <p
                    class="mt-3 rounded-lg bg-red-500/10 px-3 py-2 text-sm text-red-700"
                >
                    {formError}
                </p>
            {/if}

            <div
                class="mt-4 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"
            >
                <button type="button" class="btn-ghost" onclick={closeForm}
                    >Batal</button
                >
                <button type="submit" class="btn-primary" disabled={saving}>
                    {saving
                        ? "Menyimpan..."
                        : editingId
                          ? "Simpan Perubahan"
                          : "Tambah Soal"}
                </button>
            </div>
        </form>
    {/if}

    <div class="mb-4 grid gap-2 sm:grid-cols-[1fr_10rem]">
        <input
            class="input"
            bind:value={search}
            oninput={() => (page = 1)}
            placeholder="Cari soal, pilihan, atau jawaban..."
        />
        <select
            class="input"
            bind:value={lessonFilter}
            onchange={() => (page = 1)}
        >
            <option value="all">Semua lesson</option>
            {#each lessons as lesson}
                <option value={lesson}>{lessonLabel(lesson)}</option>
            {/each}
        </select>
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
            <p class="font-semibold text-ink">
                {search ? "Tidak ada hasil" : "Belum ada soal latihan"}
            </p>
            <p class="mt-1 text-sm text-navy/60">
                {search
                    ? "Coba kata kunci lain."
                    : 'Klik "Tambah Soal" untuk memulai.'}
            </p>
        </div>
    {:else}
        <ul class="space-y-2.5">
            {#each paged as item (item.lesson + "-" + item.id)}
                <li class="card p-3.5 sm:p-4">
                    <div class="flex items-start gap-3">
                        <div class="min-w-0 flex-1 space-y-2">
                            <div class="flex flex-wrap gap-1.5">
                                <span
                                    class="rounded-md bg-sand/40 px-2 py-0.5 text-[11px] font-medium text-navy/70"
                                    >{lessonLabel(item.lesson)}</span
                                >
                                <span
                                    class="rounded-md bg-sage/15 px-2 py-0.5 text-[11px] font-medium text-navy/70"
                                    >{item.tipe_soal}</span
                                >
                            </div>

                            <BunpouSentence
                                kalimat={item.soal}
                                showTerjemahan={false}
                                className="text-base font-semibold text-ink"
                            />

                            {#if item.tipe_soal === "pilihan ganda"}
                                <div class="flex flex-wrap gap-1.5">
                                    {#each item.pilihan_jawaban || [] as choice, index}
                                        <span
                                            class="rounded-md px-2 py-1 text-xs font-medium {isCorrectChoice(
                                                item,
                                                choice,
                                            )
                                                ? 'bg-sage/20 text-navy'
                                                : 'bg-navy/5 text-navy/60'}"
                                        >
                                            {String.fromCharCode(65 + index)}.
                                            {choice}
                                        </span>
                                    {/each}
                                </div>
                            {:else}
                                <p class="text-xs text-navy/60">
                                    Jawaban:
                                    <span class="font-jp font-semibold text-ink"
                                        >{item.jawaban_benar}</span
                                    >
                                </p>
                            {/if}
                        </div>

                        <div class="flex shrink-0 gap-1.5">
                            <button
                                class="rounded-lg border border-steel/25 px-2.5 py-1.5 text-xs font-semibold text-navy hover:bg-cream/70"
                                onclick={() => openEdit(item)}
                            >
                                Ubah
                            </button>
                            <button
                                class="rounded-lg border border-red-300 px-2.5 py-1.5 text-xs font-semibold text-red-600 hover:bg-red-500/10"
                                onclick={() => (deleteTarget = item)}
                            >
                                Hapus
                            </button>
                        </div>
                    </div>
                </li>
            {/each}
        </ul>
        <p class="mt-3 text-center text-xs text-navy/50">
            Menampilkan {rangeStart}–{rangeEnd} dari {filtered.length} soal
        </p>
        <Pagination
            page={currentPage}
            {totalPages}
            onPageChange={(p) => (page = p)}
        />
    {/if}
</div>

{#if deleteTarget}
    <div
        class="fixed inset-0 z-40 flex items-center justify-center bg-ink/60 p-4 backdrop-blur-sm"
    >
        <div class="card w-full max-w-sm p-5 sm:p-6">
            <h3 class="text-lg font-bold text-ink">Hapus soal ini?</h3>
            <p class="mt-2 text-sm text-navy/70">
                Soal pada {lessonLabel(deleteTarget.lesson)} akan dihapus.
            </p>
            <p class="mt-1 text-xs text-navy/50">
                Tindakan ini tidak dapat dibatalkan.
            </p>

            <div
                class="mt-5 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"
            >
                <button
                    class="btn-ghost"
                    onclick={() => (deleteTarget = null)}
                    disabled={deleting}>Batal</button
                >
                <button
                    class="btn-danger"
                    onclick={confirmDelete}
                    disabled={deleting}
                >
                    {deleting ? "Menghapus..." : "Hapus"}
                </button>
            </div>
        </div>
    </div>
{/if}

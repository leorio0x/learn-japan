<script>
    import { onMount } from "svelte";
    import { api } from "../lib/api.js";
    import Furigana from "../components/Furigana.svelte";
    import Pagination from "../components/Pagination.svelte";

    const bukuOptions = ["minna no nihonggo 1", "minna no nihonggo 2"];

    let words = $state([]);
    let loading = $state(true);
    let error = $state("");
    let search = $state("");

    let showForm = $state(false);
    let editingId = $state(null);
    let saving = $state(false);
    let formError = $state("");
    let form = $state(emptyForm());

    let deleteTarget = $state(null);
    let deleting = $state(false);

    let selectMode = $state(false);
    let selectedIds = $state([]);
    let showBulkDeleteConfirm = $state(false);
    let bulkDeleting = $state(false);

    let showBulk = $state(false);
    let bulkInputEl = $state(null);
    let bulkFile = $state(null);
    let bulkUploading = $state(false);
    let bulkResult = $state(null);
    let bulkError = $state("");

    const filtered = $derived(
        words.filter((w) => {
            const q = search.trim().toLowerCase();
            if (!q) return true;
            return [
                w.kanji,
                w.furigana,
                ...((w.contoh_percakapan || []).flatMap((c) => [
                    c.pertanyaan,
                    c.arti_pertanyaan,
                    c.jawaban,
                    c.arti_jawaban,
                ])),
                (w.terjemahan || []).join(" "),
                w.buku,
                w.lesson,
            ]
                .filter(Boolean)
                .some((v) => String(v).toLowerCase().includes(q));
        }),
    );

    const allSelected = $derived(
        filtered.length > 0 &&
            filtered.every((item) => selectedIds.includes(item.id)),
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

    onMount(load);

    function emptyForm() {
        return {
            kanji: "",
            kana: "",
            terjemahan: [""],
            contoh_percakapan: [],
            buku: "minna no nihonggo 1",
            lesson: "1",
        };
    }

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

    function addArti() {
        form.terjemahan = [...form.terjemahan, ""];
    }

    function removeArti(i) {
        form.terjemahan = form.terjemahan.filter((_, idx) => idx !== i);
    }

    function addPercakapan() {
        form.contoh_percakapan = [
            ...form.contoh_percakapan,
            {
                pertanyaan: "",
                arti_pertanyaan: "",
                jawaban: "",
                arti_jawaban: "",
            },
        ];
    }

    function removePercakapan(i) {
        form.contoh_percakapan = form.contoh_percakapan.filter(
            (_, idx) => idx !== i,
        );
    }

    function openCreate() {
        form = emptyForm();
        editingId = null;
        formError = "";
        showForm = true;
    }

    function openEdit(w) {
        form = {
            kanji: w.kanji || "",
            kana: w.furigana || "",
            terjemahan: (w.terjemahan || []).length ? [...w.terjemahan] : [""],
            contoh_percakapan: (w.contoh_percakapan || []).map((c) => ({
                pertanyaan: c.pertanyaan || "",
                arti_pertanyaan: c.arti_pertanyaan || "",
                jawaban: c.jawaban || "",
                arti_jawaban: c.arti_jawaban || "",
            })),
            buku: w.buku || "minna no nihonggo 1",
            lesson: w.lesson || "1",
        };
        editingId = w.id;
        formError = "";
        showForm = true;
    }

    function closeForm() {
        showForm = false;
        editingId = null;
        formError = "";
        form = emptyForm();
    }

    async function submit(event) {
        event.preventDefault();
        formError = "";
        const arti = form.terjemahan.map((a) => a.trim()).filter(Boolean);
        if (!arti.length) {
            formError = "Minimal satu terjemahan wajib diisi";
            return;
        }
        if (!form.kanji.trim() && !form.kana.trim()) {
            formError = "Isi kanji atau furigana";
            return;
        }

        const percakapan = form.contoh_percakapan
            .map((c) => ({
                pertanyaan: c.pertanyaan.trim(),
                arti_pertanyaan: c.arti_pertanyaan.trim(),
                jawaban: c.jawaban.trim(),
                arti_jawaban: c.arti_jawaban.trim(),
            }))
            .filter((c) =>
                c.pertanyaan ||
                c.arti_pertanyaan ||
                c.jawaban ||
                c.arti_jawaban
            );

        const payload = {
            kanji: form.kanji.trim(),
            furigana: form.kana.trim(),
            terjemahan: arti,
            contoh_percakapan: percakapan,
            buku: form.buku.trim(),
            lesson: form.lesson.trim(),
        };

        saving = true;
        try {
            if (editingId) {
                await api.update(editingId, payload);
            } else {
                await api.create(payload);
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
        try {
            await api.remove(deleteTarget.id);
            deleteTarget = null;
            await load();
        } catch (e) {
            error = e.message;
        } finally {
            deleting = false;
        }
    }

    function toggleSelectMode() {
        selectMode = !selectMode;
        if (!selectMode) selectedIds = [];
    }

    function toggleSelect(id) {
        selectedIds = selectedIds.includes(id)
            ? selectedIds.filter((x) => x !== id)
            : [...selectedIds, id];
    }

    function isSelected(id) {
        return selectedIds.includes(id);
    }

    function toggleSelectAll() {
        selectedIds = allSelected ? [] : filtered.map((item) => item.id);
    }

    async function confirmBulkDelete() {
        if (!selectedIds.length) return;
        bulkDeleting = true;
        error = "";
        try {
            await api.bulkRemoveKotoba(selectedIds);
            selectedIds = [];
            selectMode = false;
            showBulkDeleteConfirm = false;
            await load();
        } catch (e) {
            error = e.message;
        } finally {
            bulkDeleting = false;
        }
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
            bulkResult = await api.bulkKotoba(bulkFile);
            bulkFile = null;
            if (bulkInputEl) bulkInputEl.value = "";
            await load();
        } catch (e) {
            bulkError = e.message;
        } finally {
            bulkUploading = false;
        }
    }
</script>

<div class="mx-auto max-w-3xl">
    <div
        class="mb-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
    >
        <div>
            <h1 class="text-xl font-bold text-ink sm:text-2xl">
                Kelola Kotoba
            </h1>
            <p class="mt-1 text-sm text-navy/60">
                Tambah, ubah, dan hapus kosakata bahasa Jepang
            </p>
        </div>
        {#if !showForm}
            <div class="flex flex-col gap-2 sm:flex-row">
                {#if selectMode}
                    <button
                        class="btn-ghost w-full sm:w-auto"
                        onclick={toggleSelectMode}>Batal</button
                    >
                {:else}
                    <button
                        class="btn-ghost w-full sm:w-auto"
                        onclick={openBulk}>Bulk Input</button
                    >
                    <button
                        class="btn-ghost w-full sm:w-auto"
                        onclick={toggleSelectMode}>Bulk Hapus</button
                    >
                    <button
                        class="btn-primary w-full sm:w-auto"
                        onclick={openCreate}>+ Tambah Kotoba</button
                    >
                {/if}
            </div>
        {/if}
    </div>

    {#if showForm}
        <form class="card mb-5 p-4 sm:p-5" onsubmit={submit}>
            <div class="mb-4 flex items-center justify-between">
                <h2 class="font-bold text-ink">
                    {editingId ? "Ubah Kotoba" : "Kotoba Baru"}
                </h2>
                <button
                    type="button"
                    class="text-sm font-medium text-navy/50 hover:text-navy"
                    onclick={closeForm}>Tutup</button
                >
            </div>

            <div class="grid gap-3.5 sm:grid-cols-2">
                <div>
                    <label class="label" for="kanji">Kanji</label>
                    <input
                        id="kanji"
                        class="input font-jp"
                        bind:value={form.kanji}
                        placeholder="cth: 私"
                    />
                </div>
                <div>
                    <label class="label" for="kana">Furigana</label>
                    <input
                        id="kana"
                        class="input font-jp"
                        bind:value={form.kana}
                        placeholder="cth: わたし"
                    />
                </div>

                <div class="sm:col-span-2">
                    <span class="label">Terjemahan (bisa lebih dari satu)</span>
                    <div class="space-y-2">
                        {#each form.terjemahan as _, i}
                            <div class="flex gap-2">
                                <input
                                    class="input"
                                    bind:value={form.terjemahan[i]}
                                    placeholder="arti ke-{i + 1}"
                                />
                                <button
                                    type="button"
                                    class="shrink-0 rounded-xl border border-red-300 px-3 text-lg font-bold text-red-500 transition hover:bg-red-500/10 disabled:opacity-40"
                                    onclick={() => removeArti(i)}
                                    disabled={form.terjemahan.length <= 1}
                                    aria-label="Hapus arti"
                                >
                                    −
                                </button>
                            </div>
                        {/each}
                    </div>
                    <button
                        type="button"
                        class="mt-2 rounded-lg border border-steel/30 px-3 py-1.5 text-xs font-semibold text-navy transition hover:bg-cream/70"
                        onclick={addArti}
                    >
                        + Tambah arti
                    </button>
                </div>

                <div class="sm:col-span-2">
                    <div class="flex items-center justify-between">
                        <span class="label"
                            >Contoh Percakapan (bisa lebih dari satu)</span
                        >
                        <button
                            type="button"
                            class="rounded-lg border border-steel/30 px-3 py-1.5 text-xs font-semibold text-navy transition hover:bg-cream/70"
                            onclick={addPercakapan}
                        >
                            + Tambah percakapan
                        </button>
                    </div>
                    {#if !form.contoh_percakapan.length}
                        <p class="mt-1 text-xs text-navy/50">
                            Belum ada contoh percakapan.
                        </p>
                    {:else}
                        <div class="mt-2 space-y-3">
                            {#each form.contoh_percakapan as c, i}
                                <div
                                    class="rounded-xl border border-steel/20 p-3"
                                >
                                    <div
                                        class="mb-2 flex items-center justify-between"
                                    >
                                        <span
                                            class="text-xs font-semibold text-navy/60"
                                            >Percakapan {i + 1}</span
                                        >
                                        <button
                                            type="button"
                                            class="text-xs font-semibold text-red-500 hover:text-red-700"
                                            onclick={() => removePercakapan(i)}
                                        >
                                            Hapus
                                        </button>
                                    </div>
                                    <div class="grid gap-2.5 sm:grid-cols-2">
                                        <input
                                            class="input font-jp"
                                            bind:value={c.pertanyaan}
                                            placeholder="Pertanyaan"
                                        />
                                        <input
                                            class="input"
                                            bind:value={c.arti_pertanyaan}
                                            placeholder="Arti pertanyaan"
                                        />
                                        <input
                                            class="input font-jp"
                                            bind:value={c.jawaban}
                                            placeholder="Jawaban"
                                        />
                                        <input
                                            class="input"
                                            bind:value={c.arti_jawaban}
                                            placeholder="Arti jawaban"
                                        />
                                    </div>
                                </div>
                            {/each}
                        </div>
                    {/if}
                </div>

                <div>
                    <label class="label" for="buku">Buku</label>
                    <select id="buku" class="input" bind:value={form.buku}>
                        {#each bukuOptions as b}
                            <option value={b}>{b}</option>
                        {/each}
                    </select>
                </div>
                <div>
                    <label class="label" for="lesson">Lesson</label>
                    <input
                        id="lesson"
                        class="input"
                        bind:value={form.lesson}
                        placeholder="1"
                    />
                </div>
            </div>

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
                          : "Tambah Kotoba"}
                </button>
            </div>
        </form>
    {/if}

    <div class="mb-4">
        <input
            class="input"
            bind:value={search}
            oninput={() => (page = 1)}
            placeholder="Cari kotoba, arti, buku, atau lesson..."
        />
    </div>

    {#if selectMode}
        <div
            class="card mb-4 flex flex-col gap-2 p-3 sm:flex-row sm:items-center sm:justify-between"
        >
            <p class="text-sm font-medium text-navy/70">
                {selectedIds.length} dari {filtered.length} dipilih
            </p>
            <div class="flex flex-col gap-2 sm:flex-row">
                <button
                    class="btn-ghost w-full sm:w-auto"
                    onclick={toggleSelectAll}
                >
                    {allSelected ? "Kosongkan" : "Pilih semua"}
                </button>
                <button
                    class="btn-danger w-full sm:w-auto"
                    onclick={() => (showBulkDeleteConfirm = true)}
                    disabled={!selectedIds.length}
                >
                    Hapus Terpilih
                </button>
            </div>
        </div>
    {/if}

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
                {search ? "Tidak ada hasil" : "Belum ada kotoba"}
            </p>
            <p class="mt-1 text-sm text-navy/60">
                {search
                    ? "Coba kata kunci lain."
                    : 'Klik "Tambah Kotoba" untuk memulai.'}
            </p>
        </div>
    {:else}
        <ul class="space-y-2.5">
            {#each paged as w (w.id)}
                <li class="card flex items-center gap-3 p-3.5 sm:gap-4 sm:p-4">
                    {#if selectMode}
                        <input
                            type="checkbox"
                            class="h-4 w-4 shrink-0 accent-steel"
                            checked={isSelected(w.id)}
                            onchange={() => toggleSelect(w.id)}
                            aria-label="Pilih kotoba"
                        />
                    {/if}
                    <div class="min-w-0 flex-1">
                        <div
                            class="flex flex-wrap items-baseline gap-x-3 gap-y-1"
                        >
                            <Furigana
                                kanji={w.kanji}
                                kana={w.furigana}
                                className="font-jp text-xl font-bold text-ink"
                            />
                            <span
                                class="text-base font-semibold text-sage sm:text-lg"
                                >{(w.terjemahan || []).join(", ")}</span
                            >
                        </div>
                        <div class="mt-1.5 flex flex-wrap gap-1.5">
                            <span
                                class="rounded-md bg-navy/5 px-2 py-0.5 text-[11px] font-medium text-navy/70"
                                >{w.buku}</span
                            >
                            <span
                                class="rounded-md bg-sand/40 px-2 py-0.5 text-[11px] font-medium text-navy/70"
                                >Lesson {w.lesson}</span
                            >
                        </div>
                    </div>

                    {#if !selectMode}
                        <div class="flex shrink-0 gap-1.5">
                            <button
                                class="rounded-lg border border-steel/25 px-2.5 py-1.5 text-xs font-semibold text-navy hover:bg-cream/70"
                                onclick={() => openEdit(w)}
                            >
                                Ubah
                            </button>
                            <button
                                class="rounded-lg border border-red-300 px-2.5 py-1.5 text-xs font-semibold text-red-600 hover:bg-red-500/10"
                                onclick={() => (deleteTarget = w)}
                            >
                                Hapus
                            </button>
                        </div>
                    {/if}
                </li>
            {/each}
        </ul>
        <p class="mt-3 text-center text-xs text-navy/50">
            Menampilkan {rangeStart}–{rangeEnd} dari {filtered.length} kotoba
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
            <h3 class="text-lg font-bold text-ink">Hapus kotoba ini?</h3>
            <p class="mt-2 text-sm text-navy/70">
                <span class="font-jp font-semibold text-ink"
                    >{deleteTarget.kanji ||
                        deleteTarget.furigana}</span
                >
                &mdash; {(deleteTarget.terjemahan || []).join(", ")}
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

{#if showBulkDeleteConfirm}
    <div
        class="fixed inset-0 z-40 flex items-center justify-center bg-ink/60 p-4 backdrop-blur-sm"
    >
        <div class="card w-full max-w-sm p-5 sm:p-6">
            <h3 class="text-lg font-bold text-ink">Hapus kotoba terpilih?</h3>
            <p class="mt-2 text-sm text-navy/70">
                {selectedIds.length} kotoba akan dihapus permanen.
            </p>
            <p class="mt-1 text-xs text-navy/50">
                Tindakan ini tidak dapat dibatalkan.
            </p>

            <div
                class="mt-5 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"
            >
                <button
                    class="btn-ghost"
                    onclick={() => (showBulkDeleteConfirm = false)}
                    disabled={bulkDeleting}>Batal</button
                >
                <button
                    class="btn-danger"
                    onclick={confirmBulkDelete}
                    disabled={bulkDeleting}
                >
                    {bulkDeleting ? "Menghapus..." : "Hapus"}
                </button>
            </div>
        </div>
    </div>
{/if}

{#if showBulk}
    <div
        class="fixed inset-0 z-40 flex items-center justify-center bg-ink/60 p-4 backdrop-blur-sm"
    >
        <div class="card w-full max-w-md p-5 sm:p-6">
            <h3 class="text-lg font-bold text-ink">Bulk Input Kotoba</h3>
            <p class="mt-2 text-sm text-navy/70">
                Unggah file <span class="font-semibold">.csv</span>,
                <span class="font-semibold">.xlsx</span>, atau
                <span class="font-semibold">.json</span>. Header csv/xlsx:
            </p>
            <p
                class="mt-1.5 overflow-x-auto rounded-lg bg-navy/5 px-3 py-2 font-mono text-xs text-navy/80"
            >
                kanji, furigana, terjemahan, buku, lesson
            </p>
            <p class="mt-2 text-xs text-navy/50">
                Pemisah banyak arti: <span class="font-mono">|</span> untuk csv,
                <span class="font-mono">,</span> untuk excel. Field csv yang
                mengandung koma wajib dibungkus tanda kutip ganda, contoh:
                <span class="font-mono">"Maaf, saya mau pamit dulu."</span>.
                JSON berupa array objek sesuai format
                <span class="font-mono">data/format/kotoba.json</span>, dengan
                <span class="font-mono">terjemahan</span> berupa array dan
                <span class="font-mono">contoh_percakapan</span> opsional. ID
                dibuat otomatis.
            </p>

            <input
                type="file"
                accept=".csv,.xlsx,.json"
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
                    Berhasil mengimpor {bulkResult.imported} kotoba.
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

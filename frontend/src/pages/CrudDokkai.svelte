<script>
    import { onMount } from "svelte";
    import { api } from "../lib/api.js";
    import { stripJapaneseSpaces } from "../lib/paragraph.js";

    let { onNavigate = () => {} } = $props();

    let levels = $state([]);
    let level = $state("");
    let stories = $state([]);
    let loading = $state(true);
    let loadingStories = $state(false);
    let error = $state("");
    let search = $state("");

    let editTarget = $state(null);
    let form = $state(emptyForm());
    let saving = $state(false);
    let formError = $state("");

    let deleteTarget = $state(null);
    let deleting = $state(false);

    let showBulk = $state(false);
    let bulkInputEl = $state(null);
    let bulkFile = $state(null);
    let bulkLevel = $state("");
    let bulkUploading = $state(false);
    let bulkResult = $state(null);
    let bulkError = $state("");

    const filtered = $derived(
        stories.filter((story) => {
            const q = search.trim().toLowerCase();
            if (!q) return true;
            return [story.judul, story.paragraf, story.terjemahan, story.id]
                .filter(Boolean)
                .some((v) => String(v).toLowerCase().includes(q));
        }),
    );

    onMount(loadLevels);

    function emptyForm() {
        return { judul: "", paragraf: "", terjemahan: "" };
    }

    async function loadLevels() {
        loading = true;
        error = "";
        try {
            levels = (await api.listDokkaiLevels()) || [];
            if (!level || !levels.includes(level)) {
                level = levels.length ? levels[0] : "";
            }
            await fetchStories();
        } catch (e) {
            error = e.message;
        } finally {
            loading = false;
        }
    }

    async function fetchStories() {
        if (!level) {
            stories = [];
            return;
        }
        loadingStories = true;
        error = "";
        try {
            stories = (await api.listDokkaiLevel(level)) || [];
        } catch (e) {
            error = e.message;
        } finally {
            loadingStories = false;
        }
    }

    async function selectLevel(event) {
        level = event.target.value;
        cancelEdit();
        await fetchStories();
    }

    function preview(text) {
        const plain = stripJapaneseSpaces(
            (text || "").replace(/<([^<>]+)>/g, (_, raw) => {
                const first = raw.indexOf(":");
                return first >= 0 ? raw.slice(0, first) : raw;
            }),
        );
        return plain.length > 220 ? `${plain.slice(0, 220)}...` : plain;
    }

    function startEdit(story) {
        editTarget = story;
        form = {
            judul: story.judul || "",
            paragraf: story.paragraf || "",
            terjemahan: story.terjemahan || "",
        };
        formError = "";
    }

    function cancelEdit() {
        editTarget = null;
        form = emptyForm();
        formError = "";
    }

    async function saveEdit(event) {
        event.preventDefault();
        if (!editTarget) return;

        const paragraf = form.paragraf.trim();
        if (!paragraf) {
            formError = "Paragraf wajib diisi";
            return;
        }

        saving = true;
        formError = "";
        try {
            await api.updateDokkaiLevel(level, editTarget.id, {
                judul: form.judul.trim(),
                paragraf,
                terjemahan: form.terjemahan.trim(),
            });
            cancelEdit();
            await fetchStories();
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
            await api.removeDokkaiLevel(
                deleteTarget.level || level,
                deleteTarget.id,
            );
            deleteTarget = null;
            await fetchStories();
        } catch (e) {
            error = e.message;
        } finally {
            deleting = false;
        }
    }

    function openBulk() {
        bulkFile = null;
        bulkLevel = level || levels[0] || "";
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
            bulkResult = await api.bulkDokkaiLevel(bulkFile, bulkLevel);
            bulkFile = null;
            if (bulkInputEl) bulkInputEl.value = "";
            await loadLevels();
        } catch (e) {
            bulkError = e.message;
        } finally {
            bulkUploading = false;
        }
    }
</script>

<div class="mx-auto max-w-3xl">
    <div
        class="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"
    >
        <div>
            <button
                type="button"
                class="mb-2 inline-flex items-center gap-1 text-xs font-semibold text-steel transition hover:text-navy"
                onclick={() => onNavigate("dokkai")}
            >
                &larr; Kembali ke Dokkai
            </button>
            <h1 class="text-xl font-bold text-ink sm:text-2xl">
                Kelola Dokkai
            </h1>
            <p class="mt-1 text-sm text-navy/60">
                Tambah, ubah, dan hapus dokkai berlevel
            </p>
        </div>
        {#if !editTarget}
            <button class="btn-primary w-full sm:w-auto" onclick={openBulk}>
                Bulk Input
            </button>
        {/if}
    </div>

    {#if editTarget}
        <form class="card mb-5 p-4 sm:p-5" onsubmit={saveEdit}>
            <div class="mb-4 flex items-center justify-between">
                <h2 class="font-bold text-ink">
                    Ubah Dokkai {editTarget.id}
                </h2>
                <button
                    type="button"
                    class="text-sm font-medium text-navy/50 hover:text-navy"
                    onclick={cancelEdit}
                >
                    Tutup
                </button>
            </div>

            <div class="space-y-3.5">
                <div>
                    <label class="label" for="dokkai-judul">Judul</label>
                    <input
                        id="dokkai-judul"
                        class="input font-jp"
                        bind:value={form.judul}
                        placeholder="Judul dokkai"
                    />
                </div>
                <div>
                    <label class="label" for="dokkai-paragraf">Paragraf</label>
                    <textarea
                        id="dokkai-paragraf"
                        class="input input-textarea font-jp"
                        rows="8"
                        bind:value={form.paragraf}
                        placeholder="<kanji:kana:arti>"
                    ></textarea>
                    <p class="mt-1 text-xs text-navy/50">
                        Gunakan format
                        <span class="font-mono">&lt;kanji:kana:arti&gt;</span>
                        untuk kata yang bisa diketuk.
                    </p>
                </div>
                <div>
                    <label class="label" for="dokkai-terjemahan"
                        >Terjemahan</label
                    >
                    <textarea
                        id="dokkai-terjemahan"
                        class="input input-textarea"
                        rows="5"
                        bind:value={form.terjemahan}
                        placeholder="Terjemahan dokkai"
                    ></textarea>
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
                <button type="button" class="btn-ghost" onclick={cancelEdit}>
                    Batal
                </button>
                <button type="submit" class="btn-primary" disabled={saving}>
                    {saving ? "Menyimpan..." : "Simpan Perubahan"}
                </button>
            </div>
        </form>
    {/if}

    <div class="card mb-4 p-4 sm:p-5">
        <div class="flex flex-col gap-3 sm:flex-row sm:items-end">
            <div class="sm:w-44">
                <label class="label" for="crud-dokkai-level">Level</label>
                <select
                    id="crud-dokkai-level"
                    class="input"
                    bind:value={level}
                    onchange={selectLevel}
                    disabled={!levels.length}
                >
                    {#if !levels.length}
                        <option value="">Tidak ada level</option>
                    {/if}
                    {#each levels as l}
                        <option value={l}>{l}</option>
                    {/each}
                </select>
            </div>
            <div class="flex-1">
                <label class="label" for="dokkai-search">Cari</label>
                <input
                    id="dokkai-search"
                    class="input"
                    bind:value={search}
                    placeholder="Cari judul, isi, atau terjemahan..."
                />
            </div>
        </div>
        <p class="mt-2 text-xs text-navy/50">
            {filtered.length} dari {stories.length} dokkai
            {#if level}(level {level}){/if}
        </p>
    </div>

    {#if loading || loadingStories}
        <div class="card p-12 text-center text-navy/60">
            <span class="animate-pulse">Memuat data...</span>
        </div>
    {:else if error}
        <div class="card border-red-300 p-8 text-center text-red-700">
            <p class="font-semibold">{error}</p>
            <button class="btn-primary mt-4" onclick={loadLevels}>
                Coba lagi
            </button>
        </div>
    {:else if !filtered.length}
        <div class="card p-10 text-center">
            <p class="font-semibold text-ink">
                {search ? "Tidak ada hasil" : "Belum ada dokkai"}
            </p>
            <p class="mt-1 text-sm text-navy/60">
                {search
                    ? "Coba kata kunci lain."
                    : 'Klik "Bulk Input" untuk menambahkan dokkai.'}
            </p>
        </div>
    {:else}
        <ul class="space-y-2.5">
            {#each filtered as story (story.id)}
                <li class="card p-3.5 sm:p-4">
                    <div class="flex items-start gap-3">
                        <span
                            class="mt-0.5 shrink-0 rounded-md bg-navy/5 px-2 py-0.5 font-mono text-[11px] font-semibold text-navy/60"
                            >#{story.id}</span
                        >
                        <div class="min-w-0 flex-1">
                            <div
                                class="flex flex-wrap items-center gap-2 text-sm"
                            >
                                <span
                                    class="font-jp font-bold text-ink"
                                    >{story.judul || `Dokkai ${story.id}`}</span
                                >
                                <span
                                    class="rounded-md bg-sand/40 px-2 py-0.5 text-[11px] font-medium text-navy/70"
                                    >{story.level || level}</span
                                >
                            </div>
                            <p
                                class="mt-1 line-clamp-3 whitespace-pre-line font-jp text-sm leading-relaxed text-navy/70"
                            >
                                {preview(story.paragraf)}
                            </p>
                            {#if story.terjemahan}
                                <p
                                    class="mt-1.5 line-clamp-2 whitespace-pre-line text-xs leading-relaxed text-navy/50"
                                >
                                    {story.terjemahan}
                                </p>
                            {/if}
                        </div>
                        <div class="flex shrink-0 gap-1.5">
                            <button
                                class="rounded-lg border border-steel/25 px-2.5 py-1.5 text-xs font-semibold text-navy hover:bg-cream/70"
                                onclick={() => startEdit(story)}
                            >
                                Ubah
                            </button>
                            <button
                                class="rounded-lg border border-red-300 px-2.5 py-1.5 text-xs font-semibold text-red-600 hover:bg-red-500/10"
                                onclick={() => (deleteTarget = story)}
                            >
                                Hapus
                            </button>
                        </div>
                    </div>
                </li>
            {/each}
        </ul>
    {/if}
</div>

{#if deleteTarget}
    <div
        class="fixed inset-0 z-40 flex items-center justify-center bg-ink/60 p-4 backdrop-blur-sm"
    >
        <div class="card w-full max-w-sm p-5 sm:p-6">
            <h3 class="text-lg font-bold text-ink">Hapus dokkai ini?</h3>
            <p class="mt-2 text-sm text-navy/70">
                <span class="font-jp font-semibold text-ink"
                    >{deleteTarget.judul || `Dokkai ${deleteTarget.id}`}</span
                >
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

{#if showBulk}
    <div
        class="fixed inset-0 z-40 flex items-center justify-center bg-ink/60 p-4 backdrop-blur-sm"
    >
        <div class="card w-full max-w-md p-5 sm:p-6">
            <h3 class="text-lg font-bold text-ink">Bulk Input Dokkai</h3>
            <p class="mt-2 text-sm text-navy/70">
                Unggah file <span class="font-semibold">.json</span> berisi
                array objek sesuai format
                <span class="font-mono">data/format/dokkai-level.json</span>.
            </p>
            <p
                class="mt-1.5 overflow-x-auto rounded-lg bg-navy/5 px-3 py-2 font-mono text-xs text-navy/80"
            >
                [&#123;"id":"","level":"N5","judul":"...","paragraf":"&lt;kanji:kana:arti&gt;","terjemahan":"..."&#125;]
            </p>
            <p class="mt-2 text-xs text-navy/50">
                ID dibuat otomatis mengikuti data terakhir. Level di setiap data
                dipakai bila diisi; jika kosong memakai level cadangan di bawah.
            </p>

            <div class="mt-4">
                <label class="label" for="bulk-level"
                    >Level cadangan (opsional)</label
                >
                <select id="bulk-level" class="input" bind:value={bulkLevel}>
                    <option value="">-</option>
                    {#each levels as l}
                        <option value={l}>{l}</option>
                    {/each}
                </select>
            </div>

            <input
                type="file"
                accept=".json,application/json"
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
                    Berhasil mengimpor {bulkResult.imported} dokkai.
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

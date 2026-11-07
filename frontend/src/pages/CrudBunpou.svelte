<script>
    import { onMount } from "svelte";
    import { api } from "../lib/api.js";
    import { lessonLabel } from "../lib/bunpou.js";
    import BunpouRumus from "../components/BunpouRumus.svelte";
    import BunpouSentence from "../components/BunpouSentence.svelte";
    import Pagination from "../components/Pagination.svelte";

    const bukuOptions = ["Minna no Nihongo 1", "Minna no Nihongo 2"];

    let items = $state([]);
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

    function searchValues(b) {
        const values = [b.buku, b.lesson];
        for (const p of b.pola || []) {
            values.push(p.no_bunpou, p.nama_pola, p.rumus, p.deskripsi);
            for (const c of p.contoh || []) {
                values.push(c.pertanyaan?.kalimat, c.pertanyaan?.terjemahan);
                values.push(c.jawaban?.kalimat, c.jawaban?.terjemahan);
            }
        }
        return values;
    }

    const filtered = $derived(
        items.filter((b) => {
            const q = search.trim().toLowerCase();
            if (!q) return true;
            return searchValues(b)
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

    function emptyKalimat() {
        return { kalimat: "", terjemahan: "" };
    }

    function emptyContoh() {
        return { pertanyaan: null, jawaban: emptyKalimat() };
    }

    function emptyPola() {
        return {
            no_bunpou: "",
            nama_pola: "",
            rumus: "",
            deskripsi: "",
            contoh: [emptyContoh()],
        };
    }

    function emptyForm() {
        return {
            lesson: "",
            buku: bukuOptions[0],
            pola: [emptyPola()],
        };
    }

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

    function addPola() {
        form.pola = [...form.pola, emptyPola()];
    }
    function removePola(pi) {
        form.pola = form.pola.filter((_, i) => i !== pi);
    }
    function addContoh(pi) {
        form.pola[pi].contoh = [...form.pola[pi].contoh, emptyContoh()];
    }
    function removeContoh(pi, ci) {
        form.pola[pi].contoh = form.pola[pi].contoh.filter(
            (_, i) => i !== ci,
        );
    }
    function addPertanyaan(pi, ci) {
        form.pola[pi].contoh[ci].pertanyaan = emptyKalimat();
    }
    function removePertanyaan(pi, ci) {
        form.pola[pi].contoh[ci].pertanyaan = null;
    }

    function toFormKalimat(k) {
        return {
            kalimat: k?.kalimat || "",
            terjemahan: k?.terjemahan || "",
        };
    }

    function toFormContoh(c) {
        return {
            pertanyaan: c?.pertanyaan ? toFormKalimat(c.pertanyaan) : null,
            jawaban: toFormKalimat(c?.jawaban),
        };
    }

    function openCreate() {
        form = emptyForm();
        editingId = null;
        formError = "";
        showForm = true;
    }

    function openEdit(b) {
        form = {
            lesson: b.lesson || "",
            buku: b.buku || bukuOptions[0],
            pola: (b.pola || []).length
                ? b.pola.map((p) => ({
                      no_bunpou: p.no_bunpou || "",
                      nama_pola: p.nama_pola || "",
                      rumus: p.rumus || "",
                      deskripsi: p.deskripsi || "",
                      contoh: (p.contoh || []).length
                          ? p.contoh.map(toFormContoh)
                          : [emptyContoh()],
                  }))
                : [emptyPola()],
        };
        editingId = b.id;
        formError = "";
        showForm = true;
    }

    function closeForm() {
        showForm = false;
        editingId = null;
        formError = "";
        form = emptyForm();
    }

    function cleanKalimat(k) {
        if (!k) return null;
        const kalimat = String(k.kalimat || "").trim();
        const terjemahan = String(k.terjemahan || "").trim();
        if (!kalimat) return null;
        return { kalimat, terjemahan };
    }

    function buildPayload() {
        const pola = form.pola
            .map((p) => ({
                no_bunpou: String(p.no_bunpou || "").trim(),
                nama_pola: String(p.nama_pola || "").trim(),
                rumus: String(p.rumus || "").trim(),
                deskripsi: String(p.deskripsi || "").trim(),
                contoh: p.contoh
                    .map((c) => ({
                        pertanyaan: cleanKalimat(c.pertanyaan),
                        jawaban: cleanKalimat(c.jawaban),
                    }))
                    .filter((c) => c.pertanyaan || c.jawaban),
            }))
            .filter(
                (p) =>
                    p.no_bunpou ||
                    p.nama_pola ||
                    p.rumus ||
                    p.deskripsi ||
                    p.contoh.length,
            );
        return {
            lesson: String(form.lesson || "").trim(),
            buku: String(form.buku || "").trim(),
            pola,
        };
    }

    async function submit(event) {
        event.preventDefault();
        formError = "";

        const payload = buildPayload();

        if (!payload.pola.length) {
            formError = "Minimal satu pola wajib diisi";
            return;
        }
        for (const [pi, p] of payload.pola.entries()) {
            if (!p.nama_pola && !p.rumus) {
                formError = `Pola ke-${pi + 1} wajib memiliki nama pola atau rumus`;
                return;
            }
            if (!p.contoh.length) {
                formError = `Pola ke-${pi + 1} harus memiliki minimal satu contoh`;
                return;
            }
            for (const [ci, c] of p.contoh.entries()) {
                if (!c.jawaban) {
                    formError = `Contoh ke-${ci + 1} pada pola ke-${pi + 1} harus memiliki jawaban`;
                    return;
                }
            }
        }

        saving = true;
        try {
            if (editingId) {
                await api.updateBunpou(editingId, payload);
            } else {
                await api.createBunpou(payload);
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
            await api.removeBunpou(deleteTarget.id);
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
            await api.bulkRemoveBunpou(selectedIds);
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
</script>

<div class="mx-auto max-w-3xl">
    <div
        class="mb-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
    >
        <div>
            <h1 class="text-xl font-bold text-ink sm:text-2xl">Kelola Bunpou</h1>
            <p class="mt-1 text-sm text-navy/60">
                Tambah, ubah, dan hapus pola bunpou beserta contohnya
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
                        onclick={toggleSelectMode}>Bulk Hapus</button
                    >
                    <button
                        class="btn-primary w-full sm:w-auto"
                        onclick={openCreate}>+ Tambah Bunpou</button
                    >
                {/if}
            </div>
        {/if}
    </div>

    {#if showForm}
        <form class="card mb-5 p-4 sm:p-5" onsubmit={submit}>
            <div class="mb-4 flex items-center justify-between">
                <h2 class="font-bold text-ink">
                    {editingId ? "Ubah Bunpou" : "Bunpou Baru"}
                </h2>
                <button
                    type="button"
                    class="text-sm font-medium text-navy/50 hover:text-navy"
                    onclick={closeForm}>Tutup</button
                >
            </div>

            <div class="grid gap-3.5 sm:grid-cols-2">
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
                        placeholder="Lesson 1"
                    />
                </div>
            </div>

            <div class="mt-4 space-y-4">
                {#each form.pola as _, pi}
                    <fieldset
                        class="rounded-xl border border-navy/10 bg-cream/40 p-3.5"
                    >
                        <legend
                            class="px-1 text-xs font-semibold uppercase tracking-wide text-navy/70"
                        >
                            Pola {pi + 1}
                        </legend>

                        <div class="grid gap-2 sm:grid-cols-[7rem_1fr_auto]">
                            <input
                                class="input"
                                bind:value={form.pola[pi].no_bunpou}
                                placeholder="No. bunpou"
                            />
                            <input
                                class="input font-jp"
                                bind:value={form.pola[pi].nama_pola}
                                placeholder="nama pola, contoh: A adalah B"
                            />
                            <button
                                type="button"
                                class="shrink-0 rounded-xl border border-red-300 px-3 text-lg font-bold text-red-500 transition hover:bg-red-500/10 disabled:opacity-40"
                                onclick={() => removePola(pi)}
                                disabled={form.pola.length <= 1}
                                aria-label="Hapus pola">✕</button
                            >
                        </div>
                        <input
                            class="input mt-2 font-jp"
                            bind:value={form.pola[pi].rumus}
                            placeholder="rumus, contoh: [～:orang|kata benda] + は + [～:orang|kata benda] + です"
                        />
                        <textarea
                            class="input mt-2"
                            rows="2"
                            bind:value={form.pola[pi].deskripsi}
                            placeholder="deskripsi pola"
                        ></textarea>

                        <div class="mt-3 space-y-3">
                            {#each form.pola[pi].contoh as _, ci}
                                <div
                                    class="rounded-lg border border-steel/20 bg-white/70 p-3"
                                >
                                    <div
                                        class="mb-2 flex items-center justify-between"
                                    >
                                        <span
                                            class="text-[11px] font-semibold uppercase tracking-wide text-navy/60"
                                            >Contoh {ci + 1}</span
                                        >
                                        <button
                                            type="button"
                                            class="rounded-lg border border-red-300 px-2 py-0.5 text-xs font-bold text-red-500 transition hover:bg-red-500/10 disabled:opacity-40"
                                            onclick={() =>
                                                removeContoh(pi, ci)}
                                            disabled={form.pola[pi].contoh
                                                .length <= 1}
                                            aria-label="Hapus contoh"
                                            >✕</button
                                        >
                                    </div>

                                    {#if form.pola[pi].contoh[ci].pertanyaan}
                                        <div class="space-y-2">
                                            <div
                                                class="flex items-center justify-between"
                                            >
                                                <p
                                                    class="text-[11px] font-semibold uppercase tracking-wide text-navy/50"
                                                >
                                                    Pertanyaan (opsional)
                                                </p>
                                                <button
                                                    type="button"
                                                    class="text-[11px] font-semibold text-red-500 hover:underline"
                                                    onclick={() =>
                                                        removePertanyaan(
                                                            pi,
                                                            ci,
                                                        )}>Hapus pertanyaan</button
                                                >
                                            </div>
                                            <textarea
                                                class="input font-jp"
                                                rows="2"
                                                bind:value={form.pola[pi]
                                                    .contoh[ci].pertanyaan
                                                    .kalimat}
                                                placeholder="kalimat pertanyaan, contoh: <はい:hai:ya>、<学生:がくせい:mahasiswa> です"
                                            ></textarea>
                                            <input
                                                class="input"
                                                bind:value={form.pola[pi]
                                                    .contoh[ci].pertanyaan
                                                    .terjemahan}
                                                placeholder="terjemahan pertanyaan"
                                            />
                                        </div>
                                    {:else}
                                        <button
                                            type="button"
                                            class="rounded-lg border border-steel/30 px-3 py-1.5 text-xs font-semibold text-navy transition hover:bg-cream/70"
                                            onclick={() =>
                                                addPertanyaan(pi, ci)}
                                            >+ Tambah pertanyaan</button
                                        >
                                    {/if}

                                    <div class="mt-3 space-y-2">
                                        <p
                                            class="text-[11px] font-semibold uppercase tracking-wide text-navy/50"
                                        >
                                            Jawaban
                                        </p>
                                        <textarea
                                            class="input font-jp"
                                            rows="2"
                                            bind:value={form.pola[pi].contoh[
                                                ci
                                            ].jawaban.kalimat}
                                            placeholder="kalimat jawaban, contoh: <学生:がくせい:mahasiswa> です"
                                        ></textarea>
                                        <input
                                            class="input"
                                            bind:value={form.pola[pi].contoh[
                                                ci
                                            ].jawaban.terjemahan}
                                            placeholder="terjemahan jawaban"
                                        />
                                    </div>
                                </div>
                            {/each}
                        </div>

                        <button
                            type="button"
                            class="mt-3 rounded-lg border border-steel/30 px-3 py-1.5 text-xs font-semibold text-navy transition hover:bg-cream/70"
                            onclick={() => addContoh(pi)}>+ Tambah contoh</button
                        >
                    </fieldset>
                {/each}
            </div>

            <button
                type="button"
                class="mt-3 rounded-lg border border-steel/30 px-3 py-1.5 text-xs font-semibold text-navy transition hover:bg-cream/70"
                onclick={addPola}>+ Tambah pola</button
            >

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
                          : "Tambah Bunpou"}
                </button>
            </div>
        </form>
    {/if}

    <div class="mb-4">
        <input
            class="input"
            bind:value={search}
            oninput={() => (page = 1)}
            placeholder="Cari bunpou, rumus, terjemahan, buku, atau lesson..."
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
                {search ? "Tidak ada hasil" : "Belum ada bunpou"}
            </p>
            <p class="mt-1 text-sm text-navy/60">
                {search
                    ? "Coba kata kunci lain."
                    : 'Klik "Tambah Bunpou" untuk memulai.'}
            </p>
        </div>
    {:else}
        <ul class="space-y-2.5">
            {#each paged as b (b.id)}
                <li class="card p-3.5 sm:p-4">
                    <div class="flex items-start gap-3">
                        {#if selectMode}
                            <input
                                type="checkbox"
                                class="mt-0.5 h-4 w-4 shrink-0 accent-steel"
                                checked={isSelected(b.id)}
                                onchange={() => toggleSelect(b.id)}
                                aria-label="Pilih bunpou"
                            />
                        {/if}
                        <div class="min-w-0 flex-1 space-y-2">
                            {#each b.pola || [] as p}
                                <div class="space-y-0.5">
                                    {#if p.no_bunpou || p.nama_pola}
                                        <p
                                            class="text-sm font-semibold text-ink"
                                        >
                                            {#if p.no_bunpou}<span
                                                    class="mr-1.5 rounded bg-sage/15 px-1.5 py-0.5 text-[10px] font-bold text-navy/70"
                                                    >{p.no_bunpou}</span
                                                >{/if}{p.nama_pola}
                                        </p>
                                    {/if}
                                    {#if p.rumus}
                                        <BunpouRumus
                                            rumus={p.rumus}
                                            className="text-xs text-steel"
                                        />
                                    {/if}
                                    {#if p.deskripsi}
                                        <p class="text-[11px] text-navy/45">
                                            {p.deskripsi}
                                        </p>
                                    {/if}
                                </div>

                                {#each p.contoh || [] as c}
                                    <div
                                        class="space-y-1.5 rounded-lg border border-navy/10 bg-cream/30 p-2.5"
                                    >
                                        {#if c.pertanyaan}
                                            <div
                                                class="flex items-start gap-2"
                                            >
                                                <span
                                                    class="mt-0.5 shrink-0 rounded-md bg-steel/10 px-1.5 py-0.5 text-[10px] font-bold uppercase text-navy/60"
                                                    >P</span
                                                >
                                                <BunpouSentence
                                                    kalimat={c.pertanyaan
                                                        .kalimat}
                                                    terjemahan={c.pertanyaan
                                                        .terjemahan}
                                                    className="text-base font-semibold text-ink"
                                                />
                                            </div>
                                        {/if}
                                        {#if c.jawaban}
                                            <div
                                                class="flex items-start gap-2"
                                            >
                                                <span
                                                    class="mt-0.5 shrink-0 rounded-md bg-sage/15 px-1.5 py-0.5 text-[10px] font-bold uppercase text-navy/60"
                                                    >{c.pertanyaan ? "J" : "Contoh"}</span
                                                >
                                                <BunpouSentence
                                                    kalimat={c.jawaban.kalimat}
                                                    terjemahan={c.jawaban
                                                        .terjemahan}
                                                    className="text-base font-semibold text-sage"
                                                />
                                            </div>
                                        {/if}
                                    </div>
                                {/each}
                            {/each}

                            <div class="flex flex-wrap gap-1.5">
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
                            </div>
                        </div>

                        {#if !selectMode}
                            <div class="flex shrink-0 gap-1.5">
                                <button
                                    class="rounded-lg border border-steel/25 px-2.5 py-1.5 text-xs font-semibold text-navy hover:bg-cream/70"
                                    onclick={() => openEdit(b)}
                                >
                                    Ubah
                                </button>
                                <button
                                    class="rounded-lg border border-red-300 px-2.5 py-1.5 text-xs font-semibold text-red-600 hover:bg-red-500/10"
                                    onclick={() => (deleteTarget = b)}
                                >
                                    Hapus
                                </button>
                            </div>
                        {/if}
                    </div>
                </li>
            {/each}
        </ul>
        <p class="mt-3 text-center text-xs text-navy/50">
            Menampilkan {rangeStart}–{rangeEnd} dari {filtered.length} bunpou
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
            <h3 class="text-lg font-bold text-ink">Hapus bunpou ini?</h3>
            <p class="mt-2 text-sm text-navy/70">
                <span class="font-semibold text-ink"
                    >{deleteTarget.pola?.[0]?.nama_pola ||
                        deleteTarget.pola?.[0]?.rumus ||
                        lessonLabel(deleteTarget.lesson) ||
                        "-"}</span
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

{#if showBulkDeleteConfirm}
    <div
        class="fixed inset-0 z-40 flex items-center justify-center bg-ink/60 p-4 backdrop-blur-sm"
    >
        <div class="card w-full max-w-sm p-5 sm:p-6">
            <h3 class="text-lg font-bold text-ink">Hapus bunpou terpilih?</h3>
            <p class="mt-2 text-sm text-navy/70">
                {selectedIds.length} bunpou akan dihapus permanen.
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

<script>
    import { onMount } from "svelte";
    import { api } from "../lib/api.js";
    import Furigana from "../components/Furigana.svelte";
    import Pagination from "../components/Pagination.svelte";

    const bukuOptions = ["minna no nihonggo 1", "minna no nihonggo 2"];

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

    function searchValues(k) {
        const values = [k.buku, k.lesson];
        for (const p of k.pola || []) {
            values.push(p.pola_pertanyaan, p.pola_jawaban, p.deskripsi);
            for (const c of p.contoh || []) {
                for (const d of [
                    ...(c.pertanyaan || []),
                    ...(c.jawaban || []),
                ]) {
                    values.push(
                        d.kanji,
                        d["hiragana/katakana"],
                        d.terjemahan,
                    );
                }
            }
        }
        return values;
    }

    const filtered = $derived(
        items.filter((k) => {
            const q = search.trim().toLowerCase();
            if (!q) return true;
            return searchValues(k)
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

    function emptyDialog() {
        return { kanji: "", kana: "", terjemahan: "" };
    }

    function emptyContoh() {
        return { pertanyaan: [emptyDialog()], jawaban: [emptyDialog()] };
    }

    function emptyPola() {
        return {
            pola_pertanyaan: "",
            pola_jawaban: "",
            deskripsi: "",
            contoh: [emptyContoh()],
        };
    }

    function emptyForm() {
        return {
            pola: [emptyPola()],
            buku: "minna no nihonggo 1",
            lesson: "1",
        };
    }

    async function load() {
        loading = true;
        error = "";
        try {
            items = (await api.listKaiwa()) || [];
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
        const c = form.pola[pi].contoh[ci];
        c.pertanyaan = [...c.pertanyaan, emptyDialog()];
    }
    function removePertanyaan(pi, ci, di) {
        const c = form.pola[pi].contoh[ci];
        c.pertanyaan = c.pertanyaan.filter((_, i) => i !== di);
    }
    function addJawaban(pi, ci) {
        const c = form.pola[pi].contoh[ci];
        c.jawaban = [...c.jawaban, emptyDialog()];
    }
    function removeJawaban(pi, ci, di) {
        const c = form.pola[pi].contoh[ci];
        c.jawaban = c.jawaban.filter((_, i) => i !== di);
    }

    function toFormDialog(d) {
        return {
            kanji: d.kanji || "",
            kana: d["hiragana/katakana"] || "",
            terjemahan: d.terjemahan || "",
        };
    }

    function openCreate() {
        form = emptyForm();
        editingId = null;
        formError = "";
        showForm = true;
    }

    function openEdit(k) {
        form = {
            pola: (k.pola || []).length
                ? k.pola.map((p) => ({
                      pola_pertanyaan: p.pola_pertanyaan || "",
                      pola_jawaban: p.pola_jawaban || "",
                      deskripsi: p.deskripsi || "",
                      contoh: (p.contoh || []).length
                          ? p.contoh.map((c) => ({
                                pertanyaan: (c.pertanyaan || []).length
                                    ? c.pertanyaan.map(toFormDialog)
                                    : [emptyDialog()],
                                jawaban: (c.jawaban || []).length
                                    ? c.jawaban.map(toFormDialog)
                                    : [emptyDialog()],
                            }))
                          : [emptyContoh()],
                  }))
                : [emptyPola()],
            buku: k.buku || "minna no nihonggo 1",
            lesson: k.lesson || "1",
        };
        editingId = k.id;
        formError = "";
        showForm = true;
    }

    function closeForm() {
        showForm = false;
        editingId = null;
        formError = "";
        form = emptyForm();
    }

    function cleanDialog(d) {
        return {
            kanji: d.kanji.trim(),
            "hiragana/katakana": d.kana.trim(),
            terjemahan: d.terjemahan.trim(),
        };
    }

    function isDialogEmpty(d) {
        return d.kanji.trim() === "" && d.kana.trim() === "";
    }

    async function submit(event) {
        event.preventDefault();
        formError = "";

        const pola = form.pola
            .map((p) => ({
                pola_pertanyaan: p.pola_pertanyaan.trim(),
                pola_jawaban: p.pola_jawaban.trim(),
                deskripsi: p.deskripsi.trim(),
                contoh: p.contoh
                    .map((c) => ({
                        pertanyaan: c.pertanyaan
                            .filter((d) => !isDialogEmpty(d))
                            .map(cleanDialog),
                        jawaban: c.jawaban
                            .filter((d) => !isDialogEmpty(d))
                            .map(cleanDialog),
                    }))
                    .filter((c) => c.pertanyaan.length || c.jawaban.length),
            }))
            .filter(
                (p) =>
                    p.pola_pertanyaan ||
                    p.pola_jawaban ||
                    p.deskripsi ||
                    p.contoh.length,
            );

        if (!pola.length) {
            formError = "Minimal satu pola wajib diisi";
            return;
        }
        for (const [pi, p] of pola.entries()) {
            if (!p.contoh.length) {
                formError = `Pola ke-${pi + 1} harus memiliki minimal satu contoh`;
                return;
            }
            for (const [ci, c] of p.contoh.entries()) {
                if (!c.pertanyaan.length) {
                    formError = `Contoh ke-${ci + 1} pada pola ke-${pi + 1} harus memiliki pertanyaan`;
                    return;
                }
                if (!c.jawaban.length) {
                    formError = `Contoh ke-${ci + 1} pada pola ke-${pi + 1} harus memiliki jawaban`;
                    return;
                }
            }
        }

        const payload = {
            pola,
            buku: form.buku.trim(),
            lesson: form.lesson.trim(),
        };

        saving = true;
        try {
            if (editingId) {
                await api.updateKaiwa(editingId, payload);
            } else {
                await api.createKaiwa(payload);
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
            await api.removeKaiwa(deleteTarget.id);
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
            await api.bulkRemoveKaiwa(selectedIds);
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
            <h1 class="text-xl font-bold text-ink sm:text-2xl">Kelola Kaiwa</h1>
            <p class="mt-1 text-sm text-navy/60">
                Tambah, ubah, dan hapus percakapan beserta artinya
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
                        onclick={openCreate}>+ Tambah Kaiwa</button
                    >
                {/if}
            </div>
        {/if}
    </div>

    {#if showForm}
        <form class="card mb-5 p-4 sm:p-5" onsubmit={submit}>
            <div class="mb-4 flex items-center justify-between">
                <h2 class="font-bold text-ink">
                    {editingId ? "Ubah Kaiwa" : "Kaiwa Baru"}
                </h2>
                <button
                    type="button"
                    class="text-sm font-medium text-navy/50 hover:text-navy"
                    onclick={closeForm}>Tutup</button
                >
            </div>

            <div class="space-y-4">
                {#each form.pola as _, pi}
                    <fieldset
                        class="rounded-xl border border-navy/10 bg-cream/40 p-3.5"
                    >
                        <legend
                            class="px-1 text-xs font-semibold uppercase tracking-wide text-navy/70"
                        >
                            Pola {pi + 1}
                        </legend>

                        <div class="space-y-2">
                            <div
                                class="grid gap-2 sm:grid-cols-[1fr_1fr_auto]"
                            >
                                <input
                                    class="input font-jp"
                                    bind:value={form.pola[pi].pola_pertanyaan}
                                    placeholder="pola pertanyaan"
                                />
                                <input
                                    class="input font-jp"
                                    bind:value={form.pola[pi].pola_jawaban}
                                    placeholder="pola jawaban"
                                />
                                <button
                                    type="button"
                                    class="shrink-0 rounded-xl border border-red-300 px-3 text-lg font-bold text-red-500 transition hover:bg-red-500/10 disabled:opacity-40"
                                    onclick={() => removePola(pi)}
                                    disabled={form.pola.length <= 1}
                                    aria-label="Hapus pola">−</button
                                >
                            </div>
                            <input
                                class="input"
                                bind:value={form.pola[pi].deskripsi}
                                placeholder="deskripsi pola (opsional)"
                            />
                        </div>

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
                                            >−</button
                                        >
                                    </div>

                                    <div class="space-y-2">
                                        <p
                                            class="text-[11px] font-semibold uppercase tracking-wide text-navy/50"
                                        >
                                            Pertanyaan
                                        </p>
                                        {#each form.pola[pi].contoh[ci]
                                            .pertanyaan as _, di}
                                            <div
                                                class="grid gap-2 sm:grid-cols-[1fr_1fr_1fr_auto]"
                                            >
                                                <input
                                                    class="input font-jp"
                                                    bind:value={form.pola[pi]
                                                        .contoh[ci].pertanyaan[
                                                        di
                                                    ].kanji}
                                                    placeholder="bentuk full kanji"
                                                />
                                                <input
                                                    class="input font-jp"
                                                    bind:value={form.pola[pi]
                                                        .contoh[ci].pertanyaan[
                                                        di
                                                    ].kana}
                                                    placeholder="hiragana/katakana (opsional)"
                                                />
                                                <input
                                                    class="input"
                                                    bind:value={form.pola[pi]
                                                        .contoh[ci].pertanyaan[
                                                        di
                                                    ].terjemahan}
                                                    placeholder="terjemahan"
                                                />
                                                <button
                                                    type="button"
                                                    class="shrink-0 rounded-xl border border-red-300 px-3 text-lg font-bold text-red-500 transition hover:bg-red-500/10 disabled:opacity-40"
                                                    onclick={() =>
                                                        removePertanyaan(
                                                            pi,
                                                            ci,
                                                            di,
                                                        )}
                                                    disabled={form.pola[pi]
                                                        .contoh[ci].pertanyaan
                                                        .length <= 1}
                                                    aria-label="Hapus pertanyaan"
                                                    >−</button
                                                >
                                            </div>
                                        {/each}
                                        <button
                                            type="button"
                                            class="rounded-lg border border-steel/30 px-3 py-1.5 text-xs font-semibold text-navy transition hover:bg-cream/70"
                                            onclick={() =>
                                                addPertanyaan(pi, ci)}
                                            >+ Tambah pertanyaan</button
                                        >
                                    </div>

                                    <div class="mt-3 space-y-2">
                                        <p
                                            class="text-[11px] font-semibold uppercase tracking-wide text-navy/50"
                                        >
                                            Jawaban
                                        </p>
                                        {#each form.pola[pi].contoh[ci]
                                            .jawaban as _, di}
                                            <div
                                                class="grid gap-2 sm:grid-cols-[1fr_1fr_1fr_auto]"
                                            >
                                                <input
                                                    class="input font-jp"
                                                    bind:value={form.pola[pi]
                                                        .contoh[ci].jawaban[di]
                                                        .kanji}
                                                    placeholder="bentuk full kanji"
                                                />
                                                <input
                                                    class="input font-jp"
                                                    bind:value={form.pola[pi]
                                                        .contoh[ci].jawaban[di]
                                                        .kana}
                                                    placeholder="hiragana/katakana (opsional)"
                                                />
                                                <input
                                                    class="input"
                                                    bind:value={form.pola[pi]
                                                        .contoh[ci].jawaban[di]
                                                        .terjemahan}
                                                    placeholder="terjemahan"
                                                />
                                                <button
                                                    type="button"
                                                    class="shrink-0 rounded-xl border border-red-300 px-3 text-lg font-bold text-red-500 transition hover:bg-red-500/10 disabled:opacity-40"
                                                    onclick={() =>
                                                        removeJawaban(
                                                            pi,
                                                            ci,
                                                            di,
                                                        )}
                                                    disabled={form.pola[pi]
                                                        .contoh[ci].jawaban
                                                        .length <= 1}
                                                    aria-label="Hapus jawaban"
                                                    >−</button
                                                >
                                            </div>
                                        {/each}
                                        <button
                                            type="button"
                                            class="rounded-lg border border-steel/30 px-3 py-1.5 text-xs font-semibold text-navy transition hover:bg-cream/70"
                                            onclick={() => addJawaban(pi, ci)}
                                            >+ Tambah jawaban</button
                                        >
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

            <div class="mt-4 grid gap-3.5 sm:grid-cols-2">
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
                          : "Tambah Kaiwa"}
                </button>
            </div>
        </form>
    {/if}

    <div class="mb-4">
        <input
            class="input"
            bind:value={search}
            oninput={() => (page = 1)}
            placeholder="Cari kaiwa, terjemahan, buku, atau lesson..."
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
                {search ? "Tidak ada hasil" : "Belum ada kaiwa"}
            </p>
            <p class="mt-1 text-sm text-navy/60">
                {search
                    ? "Coba kata kunci lain."
                    : 'Klik "Tambah Kaiwa" untuk memulai.'}
            </p>
        </div>
    {:else}
        <ul class="space-y-2.5">
            {#each paged as k (k.id)}
                <li class="card p-3.5 sm:p-4">
                    <div class="flex items-start gap-3">
                        {#if selectMode}
                            <input
                                type="checkbox"
                                class="mt-0.5 h-4 w-4 shrink-0 accent-steel"
                                checked={isSelected(k.id)}
                                onchange={() => toggleSelect(k.id)}
                                aria-label="Pilih kaiwa"
                            />
                        {/if}
                        <div class="min-w-0 flex-1 space-y-2">
                            {#each k.pola || [] as p}
                                {#if p.pola_pertanyaan ||
                                    p.pola_jawaban ||
                                    p.deskripsi}
                                    <div class="space-y-0.5">
                                        <p class="text-[11px] text-navy/50">
                                            <span class="font-jp"
                                                >{p.pola_pertanyaan}</span
                                            >
                                            <span class="text-navy/30">→</span>
                                            <span class="font-jp"
                                                >{p.pola_jawaban}</span
                                            >
                                        </p>
                                        {#if p.deskripsi}
                                            <p class="text-[11px] text-navy/40">
                                                {p.deskripsi}
                                            </p>
                                        {/if}
                                    </div>
                                {/if}

                                {#each p.contoh || [] as c}
                                    <div
                                        class="space-y-1.5 rounded-lg border border-navy/10 bg-cream/30 p-2.5"
                                    >
                                        {#each c.pertanyaan || [] as q}
                                            <div
                                                class="flex items-start gap-2"
                                            >
                                                <span
                                                    class="mt-0.5 shrink-0 rounded-md bg-steel/10 px-1.5 py-0.5 text-[10px] font-bold uppercase text-navy/60"
                                                    >P</span
                                                >
                                                <div class="min-w-0">
                                                    <Furigana
                                                        kanji={q.kanji}
                                                        kana={q[
                                                            "hiragana/katakana"
                                                        ]}
                                                        className="font-jp text-base font-semibold text-ink"
                                                    />
                                                    {#if q.terjemahan}
                                                        <p
                                                            class="text-xs text-navy/50"
                                                        >
                                                            {q.terjemahan}
                                                        </p>
                                                    {/if}
                                                </div>
                                            </div>
                                        {/each}

                                        {#each c.jawaban || [] as j}
                                            <div
                                                class="flex items-start gap-2"
                                            >
                                                <span
                                                    class="mt-0.5 shrink-0 rounded-md bg-sage/15 px-1.5 py-0.5 text-[10px] font-bold uppercase text-navy/60"
                                                    >J</span
                                                >
                                                <div class="min-w-0">
                                                    <Furigana
                                                        kanji={j.kanji}
                                                        kana={j[
                                                            "hiragana/katakana"
                                                        ]}
                                                        className="font-jp text-base font-semibold text-sage"
                                                    />
                                                    {#if j.terjemahan}
                                                        <p
                                                            class="text-xs text-navy/50"
                                                        >
                                                            {j.terjemahan}
                                                        </p>
                                                    {/if}
                                                </div>
                                            </div>
                                        {/each}
                                    </div>
                                {/each}
                            {/each}

                            <div class="flex flex-wrap gap-1.5">
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
                        </div>

                        {#if !selectMode}
                            <div class="flex shrink-0 gap-1.5">
                                <button
                                    class="rounded-lg border border-steel/25 px-2.5 py-1.5 text-xs font-semibold text-navy hover:bg-cream/70"
                                    onclick={() => openEdit(k)}
                                >
                                    Ubah
                                </button>
                                <button
                                    class="rounded-lg border border-red-300 px-2.5 py-1.5 text-xs font-semibold text-red-600 hover:bg-red-500/10"
                                    onclick={() => (deleteTarget = k)}
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
            Menampilkan {rangeStart}–{rangeEnd} dari {filtered.length} kaiwa
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
            <h3 class="text-lg font-bold text-ink">Hapus kaiwa ini?</h3>
            <p class="mt-2 text-sm text-navy/70">
                <span class="font-jp font-semibold text-ink"
                    >{deleteTarget.pola?.[0]?.contoh?.[0]?.pertanyaan?.[0]?.[
                        "hiragana/katakana"
                    ] ||
                        deleteTarget.pola?.[0]?.contoh?.[0]?.pertanyaan?.[0]
                            ?.kanji ||
                        "—"}</span
                >
                <span class="text-navy/40">→</span>
                <span class="font-jp font-semibold text-ink"
                    >{deleteTarget.pola?.[0]?.contoh?.[0]?.jawaban?.[0]?.[
                        "hiragana/katakana"
                    ] ||
                        deleteTarget.pola?.[0]?.contoh?.[0]?.jawaban?.[0]
                            ?.kanji ||
                        "—"}</span
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
            <h3 class="text-lg font-bold text-ink">Hapus kaiwa terpilih?</h3>
            <p class="mt-2 text-sm text-navy/70">
                {selectedIds.length} kaiwa akan dihapus permanen.
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

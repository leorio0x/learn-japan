<script>
    import { onMount } from "svelte";
    import { auth } from "../lib/auth.svelte.js";
    import FeatureIcon from "../components/FeatureIcon.svelte";

    let { onNavigate = () => {} } = $props();

    let password = $state("");
    let submitting = $state(false);
    let error = $state("");

    const links = [
        {
            id: "crud-kotoba",
            label: "Kelola Kotoba",
            desc: "Tambah, ubah, hapus kosakata",
            icon: "kotoba",
        },
        {
            id: "crud-kanji",
            label: "Kelola Kanji",
            desc: "Tambah, ubah, hapus kanji",
            icon: "kanji",
        },
        {
            id: "crud-kaiwa",
            label: "Kelola Kaiwa",
            desc: "Tambah, ubah, hapus percakapan",
            icon: "kaiwa",
        },
    ];

    onMount(() => {
        if (auth.checking) auth.refresh();
    });

    async function submit(event) {
        event.preventDefault();
        if (!password) {
            error = "Masukkan password terlebih dahulu";
            return;
        }
        submitting = true;
        error = "";
        try {
            await auth.login(password);
            password = "";
        } catch (e) {
            error = e.message;
        } finally {
            submitting = false;
        }
    }

    async function lock() {
        await auth.logout();
    }
</script>

<div class="mx-auto max-w-lg">
    {#if auth.checking}
        <div class="card p-12 text-center text-navy/60">
            <span class="animate-pulse">Memeriksa akses...</span>
        </div>
    {:else if auth.authenticated}
        <div class="card p-5 sm:p-6">
            <div class="flex items-start justify-between gap-3">
                <div class="flex items-center gap-2.5">
                    <span
                        class="flex h-9 w-9 items-center justify-center rounded-xl bg-sage text-lg font-bold text-ink"
                        >管</span
                    >
                    <div>
                        <h1 class="text-lg font-bold text-ink">Admin</h1>
                        <p class="text-xs text-navy/60">
                            Pilih halaman kelola data
                        </p>
                    </div>
                </div>
                <button
                    type="button"
                    class="rounded-lg border border-steel/25 bg-white px-3 py-1.5 text-xs font-semibold text-navy transition hover:bg-cream/70"
                    onclick={lock}
                >
                    Keluar
                </button>
            </div>

            <div class="mt-5 space-y-2.5">
                {#each links as link}
                    <button
                        type="button"
                        class="flex w-full items-center gap-3 rounded-xl border border-steel/25 bg-white px-4 py-3 text-left transition hover:border-steel/50 hover:bg-cream/60"
                        onclick={() => onNavigate(link.id)}
                    >
                        <span
                            class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-navy/5 text-steel"
                            ><FeatureIcon name={link.icon} class="h-4 w-4" /></span
                        >
                        <span class="min-w-0">
                            <span
                                class="block text-sm font-semibold text-ink"
                                >{link.label}</span
                            >
                            <span class="block text-xs text-navy/55"
                                >{link.desc}</span
                            >
                        </span>
                        <span class="ml-auto text-navy/30">›</span>
                    </button>
                {/each}
            </div>
        </div>
    {:else}
        <form class="card p-5 sm:p-6" onsubmit={submit}>
            <div class="flex items-center gap-2.5">
                <span
                    class="flex h-9 w-9 items-center justify-center rounded-xl bg-sage text-lg font-bold text-ink"
                    >鍵</span
                >
                <div>
                    <h1 class="text-lg font-bold text-ink">Akses Terbatas</h1>
                    <p class="text-xs text-navy/60">
                        Fitur kelola data memerlukan password
                    </p>
                </div>
            </div>

            <div class="mt-5">
                <label class="label" for="admin-password">Password</label>
                <input
                    id="admin-password"
                    type="password"
                    class="input"
                    autocomplete="off"
                    bind:value={password}
                    placeholder="••••••••"
                />
            </div>

            {#if error}
                <p
                    class="mt-3 rounded-lg bg-red-500/10 px-3 py-2 text-sm text-red-700"
                >
                    {error}
                </p>
            {/if}

            <button
                type="submit"
                class="btn-primary mt-4 w-full"
                disabled={submitting}
            >
                {submitting ? "Memeriksa..." : "Masuk"}
            </button>
        </form>
    {/if}
</div>

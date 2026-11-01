<script>
    import { onMount } from "svelte";
    import { auth } from "../lib/auth.svelte.js";

    let { children } = $props();

    let password = $state("");
    let submitting = $state(false);
    let error = $state("");

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

{#if auth.checking}
    <div class="card p-12 text-center text-navy/60">
        <span class="animate-pulse">Memeriksa akses...</span>
    </div>
{:else if auth.authenticated}
    <div class="mb-3 flex justify-end">
        <button
            type="button"
            class="rounded-lg border border-steel/25 bg-white px-3 py-1.5 text-xs font-semibold text-navy transition hover:bg-cream/70"
            onclick={lock}
        >
            Keluar
        </button>
    </div>
    {@render children()}
{:else}
    <div class="mx-auto max-w-sm">
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
                <label class="label" for="crud-password">Password</label>
                <input
                    id="crud-password"
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
    </div>
{/if}

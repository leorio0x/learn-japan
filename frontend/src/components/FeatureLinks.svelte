<script>
    import { auth } from "../lib/auth.svelte.js";
    import FeatureIcon from "./FeatureIcon.svelte";

    let { page = "kotoba", onNavigate = () => {} } = $props();

    const configs = {
        kotoba: [
            { id: "dokkai-kotoba", label: "Dokkai", icon: "book" },
            {
                id: "quiz-kotoba",
                label: "Tebak Kotoba",
                icon: "quiz",
                primary: true,
            },
            {
                id: "crud-kotoba",
                label: "Kelola Kotoba",
                icon: "crud",
                authRequired: true,
            },
        ],
        kanji: [
            { id: "dokkai-kanji", label: "Dokkai", icon: "book" },
            {
                id: "quiz-kanji",
                label: "Tebak Kanji",
                icon: "quiz",
                primary: true,
            },
            {
                id: "crud-kanji",
                label: "Kelola Kanji",
                icon: "crud",
                authRequired: true,
            },
        ],
        kaiwa: [
            {
                id: "quiz-kaiwa",
                label: "Tebak Kaiwa",
                icon: "quiz",
                primary: true,
            },
            {
                id: "crud-kaiwa",
                label: "Kelola Kaiwa",
                icon: "crud",
                authRequired: true,
            },
        ],
        dokkai: [
            {
                id: "crud-dokkai",
                label: "Kelola Dokkai",
                icon: "crud",
                authRequired: true,
            },
        ],
    };

    const links = $derived(
        (configs[page] || []).filter(
            (link) => !link.authRequired || auth.authenticated,
        ),
    );
</script>

{#if links.length}
    <div class="flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:justify-end">
        {#each links as link (link.id)}
            <button
                type="button"
                class="w-full justify-center sm:w-auto {link.primary
                    ? 'btn-primary'
                    : 'btn-ghost'}"
                onclick={() => onNavigate(link.id)}
            >
                <span
                    class="flex h-6 w-6 shrink-0 items-center justify-center rounded-lg {link.primary
                        ? 'bg-cream/15 text-cream'
                        : 'bg-sage/15 text-sage'}"
                >
                    <FeatureIcon name={link.icon} class="h-3.5 w-3.5" />
                </span>
                <span>{link.label}</span>
            </button>
        {/each}
    </div>
{/if}

<script>
    let { name = "quiz", class: className = "h-4 w-4" } = $props();
</script>

<svg
    class={className}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
>
    {#if name === "book"}
        <path d="M12 7v14" />
        <path
            d="M3 18a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1h5a4 4 0 0 1 4 4 4 4 0 0 1 4-4h5a1 1 0 0 1 1 1v13a1 1 0 0 1-1 1h-6a3 3 0 0 0-3 3 3 3 0 0 0-3-3z"
        />
    {:else if name === "crud"}
        <path d="M12 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
        <path
            d="M18.375 2.625a2.121 2.121 0 1 1 3 3L12 15l-4 1 1-4z"
        />
    {:else if name === "kaiwa"}
        <path
            d="M7.9 20A9 9 0 1 0 4 16.1L2 22z"
        />
        <path d="M8 12h.01" />
        <path d="M12 12h.01" />
        <path d="M16 12h.01" />
    {:else if name === "kanji"}
        <path d="M4 7V4h16v3" />
        <path d="M9 20h6" />
        <path d="M12 4v16" />
    {:else if name === "kotoba"}
        <path d="m5 8 6 6" />
        <path d="m4 14 6-6-6-6" />
        <path d="M2 5h12" />
        <path d="M7 2h1" />
        <path d="m22 22-5-10-5 10" />
        <path d="M14 18h6" />
    {:else}
        <circle cx="12" cy="12" r="10" />
        <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3" />
        <path d="M12 17h.01" />
    {/if}
</svg>

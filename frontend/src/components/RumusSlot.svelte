<script>
    import {
        computePosition,
        offset,
        flip,
        shift,
        autoUpdate,
    } from "@floating-ui/dom";

    let { value = "～", categories = [], tone = "default" } = $props();

    let open = $state(false);
    let reference = $state(null);
    let floating = $state(null);

    $effect(() => {
        if (!open || !reference || !floating) return;
        return autoUpdate(reference, floating, update);
    });

    $effect(() => {
        if (!open) return;
        const dismiss = (event) => {
            const target = event.target;
            if (reference?.contains(target)) return;
            if (floating?.contains(target)) return;
            open = false;
        };
        document.addEventListener("pointerdown", dismiss, true);
        return () =>
            document.removeEventListener("pointerdown", dismiss, true);
    });

    function update() {
        if (!reference || !floating) return;
        computePosition(reference, floating, {
            placement: "top",
            strategy: "fixed",
            middleware: [offset(8), flip(), shift({ padding: 8 })],
        }).then(({ x, y }) => {
            if (!floating) return;
            floating.style.left = `${x}px`;
            floating.style.top = `${y}px`;
        });
    }

    function show() {
        open = true;
    }

    function toggle() {
        open = !open;
    }
</script>

<button
    type="button"
    bind:this={reference}
    class="mx-0.5 inline-flex items-center rounded-md border-b-2 border-dotted px-0.5 font-jp font-semibold transition focus:outline-none focus-visible:ring-2 {tone ===
    'dark'
        ? 'border-sand/60 text-sand hover:bg-cream/10 focus-visible:ring-sand/40'
        : 'border-steel/70 text-steel hover:bg-sand/40 focus-visible:ring-steel/40'} {open
        ? tone === 'dark'
            ? 'bg-cream/10'
            : 'bg-sand/60'
        : ''}"
    onpointerenter={(event) => {
        if (event.pointerType === "mouse") show();
    }}
    onpointerleave={(event) => {
        if (event.pointerType === "mouse") open = false;
    }}
    onfocus={show}
    onclick={toggle}
    onkeydown={(event) => {
        if (event.key === "Escape") open = false;
    }}
    aria-expanded={open}
    aria-label="Lihat kategori pengisi"
>
    {value}
</button>

{#if open}
    <span
        bind:this={floating}
        class="fixed left-0 top-0 z-50 block max-w-[16rem] rounded-xl bg-ink px-3 py-2.5 text-left text-xs text-cream shadow-xl ring-1 ring-cream/10"
        role="tooltip"
    >
        <span class="block text-[10px] font-semibold uppercase tracking-wide text-cream/50"
            >Dapat diisi</span
        >
        {#if categories.length}
            <span class="mt-1.5 flex flex-wrap gap-1.5">
                {#each categories as category}
                    <span
                        class="rounded-md bg-cream/10 px-2 py-0.5 text-[11px] font-medium text-cream/90"
                        >{category}</span
                    >
                {/each}
            </span>
        {:else}
            <span class="mt-1.5 block text-cream/70"
                >Kata yang bisa diganti</span
            >
        {/if}
    </span>
{/if}

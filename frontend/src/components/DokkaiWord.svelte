<script>
    import {
        computePosition,
        offset,
        flip,
        shift,
        autoUpdate,
    } from "@floating-ui/dom";

    let { kanji = "", kana = "", arti = "", hideKanji = false } = $props();

    let open = $state(false);
    let reference = $state(null);
    let floating = $state(null);

    const label = $derived(hideKanji ? kana || kanji : kanji || kana);

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

    function toggle(event) {
        event.stopPropagation();
        open = !open;
    }

    function handleClick(event) {
        event.stopPropagation();
        toggle(event);
    }
</script>

<span
    bind:this={reference}
    class="-mx-0.5 cursor-help touch-manipulation select-none rounded px-0.5 transition hover:bg-sand/40 {open
        ? 'bg-sand/60'
        : ''}"
    role="button"
    tabindex="0"
    onclick={handleClick}
    onkeydown={(e) => {
        if (e.key === "Enter") {
            e.preventDefault();
            open = !open;
        }
    }}
    aria-label="Klik untuk melihat furigana dan arti"
>
    {label}
</span>

{#if open}
    <span
        bind:this={floating}
        class="fixed left-0 top-0 z-50 block max-w-[16rem] rounded-xl bg-ink px-3 py-2 text-left text-xs text-cream shadow-xl ring-1 ring-cream/10"
    >
        {#if kana}
            <span class="block font-jp text-sm font-semibold text-sage"
                >{kana}</span
            >
        {/if}
        {#if arti}
            <span class="mt-0.5 block text-cream/85">{arti}</span>
        {/if}
        {#if !kana && !arti}
            <span class="block text-cream/70">Tidak ada keterangan</span>
        {/if}
    </span>
{/if}

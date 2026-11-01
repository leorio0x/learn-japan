<script>
  import { onMount } from "svelte";

  let { page = "kotoba", onNavigate = () => {} } = $props();

  const items = [
    { id: "kotoba", label: "Kotoba", icon: "語" },
    { id: "kanji", label: "Kanji", icon: "漢" },
    { id: "kaiwa", label: "Kaiwa", icon: "会" },
    { id: "dokkai", label: "Dokkai", icon: "読" },
  ];

  let manualHide = $state(false);
  let scrollHide = $state(false);
  const hidden = $derived(manualHide || scrollHide);

  onMount(() => {
    let lastY = window.scrollY;
    const onScroll = () => {
      const y = window.scrollY;
      if (y <= 80) scrollHide = false;
      else if (y > lastY + 6) scrollHide = true;
      else if (y < lastY - 6) scrollHide = false;
      lastY = y;
    };
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  });
</script>

<header
  class="sticky top-0 z-30 border-b border-cream/10 bg-ink/95 backdrop-blur transition-transform duration-300 will-change-transform {hidden
    ? '-translate-y-full'
    : 'translate-y-0'}"
>
  <nav
    class="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-3 px-4 py-3 sm:px-6"
  >
    <button
      class="flex items-center gap-2.5 rounded-xl text-left"
      onclick={() => onNavigate("kotoba")}
    >
      <span
        class="flex h-9 w-9 items-center justify-center rounded-xl bg-sage text-lg font-bold text-ink"
        >日</span
      >
      <span class="leading-tight">
        <span
          class="block text-sm font-bold tracking-wide text-cream sm:text-base"
          >Learn Japan</span
        >
        <span class="hidden text-[11px] text-cream/60 sm:block"
          >Quiz Kotoba, Kanji &amp; Kaiwa</span
        >
      </span>
    </button>

    <div
      class="order-last flex w-full items-center gap-2 sm:order-none sm:w-auto"
    >
      <div
        class="grid flex-1 grid-cols-4 gap-1 rounded-xl bg-white/5 p-1 sm:flex sm:flex-none"
      >
        {#each items as item}
          <button
            class="flex items-center justify-center gap-1.5 rounded-lg px-3 py-2 text-xs font-semibold transition sm:flex-none sm:px-3.5 sm:text-sm {page ===
            item.id
              ? 'bg-steel text-cream shadow'
              : 'text-cream/70 hover:bg-white/10 hover:text-cream'}"
            onclick={() => onNavigate(item.id)}
          >
            <span class="font-jp">{item.icon}</span>
            <span>{item.label}</span>
          </button>
        {/each}
      </div>

      <button
        type="button"
        class="flex shrink-0 items-center justify-center rounded-xl bg-white/5 px-2.5 py-3 text-cream/70 transition hover:bg-white/10 hover:text-cream"
        onclick={() => (manualHide = true)}
        aria-label="Sembunyikan navbar"
        title="Sembunyikan navbar"
      >
        <span class="text-base leading-none">⌃</span>
      </button>
    </div>
  </nav>
</header>

{#if manualHide}
  <button
    type="button"
    class="fixed right-3 top-3 z-40 flex items-center gap-1.5 rounded-full bg-ink/90 px-3.5 py-2 text-xs font-semibold text-cream shadow-lg ring-1 ring-cream/10 transition hover:bg-ink"
    onclick={() => (manualHide = false)}
  >
    <span class="font-jp text-sage">日</span> Tampilkan menu
  </button>
{/if}

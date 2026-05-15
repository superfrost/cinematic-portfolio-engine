<script lang="ts">
    import { onMount, tick } from 'svelte';
    import PortfolioBlock from '$lib/components/PortfolioBlock.svelte';
    import ReviewsSection from '$lib/components/ReviewsSection.svelte';
    import AboutBlock from '$lib/components/AboutBlock.svelte';

    let { data } = $props();

    let searchQuery = $state("");

    let filteredCases = $derived(
        data.cases.filter(c =>
            c.title.toLowerCase().includes(searchQuery.toLowerCase())
        )
    );

    onMount(async () => {
        if (window.location.hash) {
            await tick();
            const id = window.location.hash.replace('#', '');
            const target = document.getElementById(id);
            if (target) {
                target.scrollIntoView({ behavior: 'smooth', block: 'start' });
            }
        }
    });
</script>

<div class="min-h-screen bg-neutral-950 text-neutral-100 antialiased font-sans selection:bg-amber-500 selection:text-neutral-950">

    <header class="fixed top-0 left-0 w-full z-50 bg-neutral-950/80 backdrop-blur-md border-b border-neutral-900">
        <div class="container mx-auto px-4 h-16 flex items-center justify-between gap-4">
            <span class="font-bold tracking-tighter uppercase text-[10px] md:text-sm flex-shrink-0">
                Видеограф | продюсер
            </span>

            <div class="relative w-full max-w-[120px] xs:max-w-[180px] md:max-w-md">
                <label for="search-input" class="sr-only">Поиск проектов</label>
                <input
                    id="search-input"
                    type="text"
                    bind:value={searchQuery}
                    placeholder="Поиск..."
                    class="w-full bg-neutral-900 border border-neutral-800 rounded-md py-1.5 px-2 text-sm focus:outline-none focus:border-amber-500 transition-colors"
                />
            </div>

            <nav class="flex gap-6 text-sm tracking-wide text-neutral-400">
                <a href="#about" class="hover:text-amber-500 transition-colors">Об авторе</a>
                <a href="#reviews" class="hover:text-amber-500 transition-colors">Отзывы</a>
                <a href="#projects" class="hover:text-amber-500 transition-colors">Работы</a>
            </nav>
        </div>
    </header>

    <section id="about" class="pt-24">
        <AboutBlock author={data.author} />
    </section>

    <ReviewsSection reviews={data.reviews} />

    <section id="projects">
        {#each filteredCases as project (project.id)}
            <PortfolioBlock {project} />
        {:else}
            <div class="py-20 text-center text-neutral-500">Проекты не найдены</div>
        {/each}
    </section>
</div>
<script lang="ts">
    import { resolve } from '$app/paths';
    import PortfolioBlock from '$lib/components/PortfolioBlock.svelte';
    import ReviewBlock from '$lib/components/ReviewBlock.svelte';

    let { data } = $props();

    // Состояние поиска
    let searchQuery = $state("");

    // Фильтрация кейсов по заголовку
    let filteredCases = $derived(
        data.cases.filter(c => 
            c.title.toLowerCase().includes(searchQuery.toLowerCase())
        )
    );

    $effect(() => {
        if (window.location.hash) {
            setTimeout(() => {
                const id = window.location.hash.replace('#', '');
                const target = document.getElementById(id);
                if (target) {
                    target.scrollIntoView({ behavior: 'smooth', block: 'start' });
                }
            }, 100);
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
                <input 
                    type="text" 
                    bind:value={searchQuery}
                    placeholder="Поиск..."
                    class="w-full bg-neutral-900 border border-neutral-800 rounded-full py-1.5 px-10 text-sm focus:outline-none focus:border-amber-500 transition-colors"
                />
                <span class="absolute left-2.5 top-1/2 -translate-y-1/2 text-[10px] md:text-xs text-neutral-500">
                </span>
            </div>

            <nav class="flex gap-6 text-sm tracking-wide text-neutral-400">
                <a href="#projects" class="hover:text-amber-500 transition-colors">Работы</a>
                <a href="#reviews" class="hover:text-amber-500 transition-colors">Отзывы</a>
                <a href="#about" class="hover:text-amber-500 transition-colors">Об авторе</a>
            </nav>
        </div>
    </header>

    <main id="projects" class="pt-24">
        {#each filteredCases as project (project.id)}
            <PortfolioBlock {project} />
        {:else}
            <div class="py-20 text-center text-neutral-500">Проекты не найдены</div>
        {/each}
    </main>

    <section id="reviews" class="w-full py-24 bg-neutral-950 border-t border-neutral-900 scroll-mt-16">
        <div class="container mx-auto px-4 max-w-6xl">
            <h2 class="text-4xl font-bold text-neutral-100 mb-12">Отзывы</h2>
            
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                {#each data.reviews as review}
                    <ReviewBlock {review} />
                {/each}
            </div>
        </div>
    </section>

    <footer id="about" class="w-full min-h-[80vh] flex items-center bg-neutral-900/20 py-20 border-t border-neutral-900 scroll-mt-16">
        <div class="container mx-auto px-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-12 items-start max-w-5xl mx-auto">
                
                <!-- ЛЕВАЯ КОЛОНКА: Медиа (Фото/Видео) -->
                <div class="flex flex-col gap-6 w-full">
                    
                    <!-- Условие на ФОТО -->
                    {#if data.author.media?.photo}
                        <div class="w-full bg-neutral-800 rounded-2xl overflow-hidden shadow-2xl border border-neutral-800 group">
                            <img 
                            src={resolve(data.author.media.photo as `/#${string}`)} 
                                alt={data.author.name} 
                                class="w-full h-auto grayscale contract-125 filter group-hover:grayscale-0 transition-all duration-700"
                            />
                        </div>
                    {/if}

                    <!-- Условие на ВИДЕО -->
                    {#if data.author.media.video}
                        <div class="aspect-video bg-neutral-950 rounded-2xl overflow-hidden border border-neutral-800 shadow-lg w-full">
                            <iframe 
                                class="w-full h-full border-0" 
                                src={data.author.media.video.provider === 'youtube' 
                                    ? `https://www.youtube.com/embed/${data.author.media.video.id}` 
                                    : `https://rutube.ru/play/embed/${data.author.media.video.id}`} 
                                title="Author showreel" 
                                allowfullscreen
                            ></iframe>
                        </div>
                    {/if}

                    {#if !data.author.media.photo && !data.author.media.video}
                        <div class="p-12 border-2 border-dashed border-neutral-800 rounded-2xl text-neutral-600 text-center">
                            Медиа-контент не указан
                        </div>
                    {/if}
                </div>

                <!-- ПРАВАЯ КОЛОНКА: Динамические данные автора -->
                <div class="flex flex-col justify-center sticky top-24">
                    <span class="text-xs uppercase tracking-widest text-amber-500 font-semibold mb-3 block">
                        {data.author.profession}
                    </span>
                    <h2 class="text-3xl md:text-5xl font-bold tracking-tight text-neutral-100 mb-6">
                        {data.author.name}
                    </h2>
                    <p class="text-neutral-400 font-light text-base md:text-lg leading-relaxed mb-8">
                        {data.author.bio}
                    </p>

                    <div class="space-y-4 border-t border-neutral-800 pt-6 font-light">
                        {#if data.author.contacts.telegram}
                            <div class="flex justify-between py-2 border-b border-neutral-900">
                                <span class="text-neutral-500 text-sm">Telegram:</span>
                                <a href="https://t.me/{data.author.contacts.telegram.replace('@', '')}" class="hover:text-amber-500 transition-colors">
                                    {data.author.contacts.telegram}
                                </a>
                            </div>
                        {/if}
                        
                        {#if data.author.contacts.email}
                            <div class="flex justify-between py-2 border-b border-neutral-900">
                                <span class="text-neutral-500 text-sm">Email:</span>
                                <a href="mailto:{data.author.contacts.email}" class="hover:text-amber-500 transition-colors">
                                    {data.author.contacts.email}
                                </a>
                            </div>
                        {/if}

                        <div class="flex justify-between py-2 border-b border-neutral-900">
                            <span class="text-neutral-500 text-sm">Локация:</span>
                            <span class="text-neutral-300">{data.author.location}</span>
                        </div>
                    </div>
                </div>

            </div>
        </div>
    </footer>
</div>
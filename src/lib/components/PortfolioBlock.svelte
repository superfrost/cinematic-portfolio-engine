<script lang="ts">
    import CarouselItem from './CarouselItem.svelte';

    interface MediaItem {
        type: 'video' | 'image' | 'text';
        provider?: 'youtube' | 'vimeo' | 'rutube';
        id?: string;
        src?: string;
        content?: string;
    }

    interface Project {
        id: string;
        title: string;
        description: string;
        items: MediaItem[];
    }

    let { project }: { project: Project } = $props();

    let sliderElement = $state<HTMLDivElement | undefined>(undefined);
    let currentIndex = $state(0);

    let showCopyTooltip = $state(false);

    let autoplayInterval: ReturnType<typeof setInterval> | undefined;
    let autoplayDestroyed = false;

    // Функция полной остановки автоплея навсегда
    function stopAutoplayForever() {
        if (autoplayDestroyed) return;
        autoplayDestroyed = true;
        if (autoplayInterval) {
            clearInterval(autoplayInterval);
            autoplayInterval = undefined;
        }
    }

    function copyProjectLink() {
        // Формируем URL с якорем (например, site.com/#project-id)
        const url = `${window.location.origin}${window.location.pathname}#${project.id}`;
        
        navigator.clipboard.writeText(url).then(() => {
            showCopyTooltip = true;
            setTimeout(() => {
                showCopyTooltip = false;
            }, 2000);
        });
    }

    function scrollToIndex(index: number): void {
        stopAutoplayForever(); // Клик по точкам пагинации останавливает автоплей
        if (!sliderElement) return;
        const width = sliderElement.clientWidth;
        sliderElement.scrollTo({
            left: width * index,
            behavior: 'smooth'
        });
        currentIndex = index;
    }

    function nextSlide(): void {
        const next = (currentIndex + 1) % project.items.length;
        if (!sliderElement) return;
        sliderElement.scrollTo({ left: sliderElement.clientWidth * next, behavior: 'smooth' });
        currentIndex = next;
    }

    function prevSlide(): void {
        stopAutoplayForever(); // Клик по стрелке назад останавливает автоплей
        const prev = (currentIndex - 1 + project.items.length) % project.items.length;
        if (!sliderElement) return;
        sliderElement.scrollTo({ left: sliderElement.clientWidth * prev, behavior: 'smooth' });
        currentIndex = prev;
    }

    function handleScroll(e: Event): void {
        // ЛЮБОЕ движение ленты (жест пальцем, скролл мышкой) убивает автоплей
        stopAutoplayForever(); 
        
        const target = e.target as HTMLDivElement;
        if (!target) return;
        
        const width = target.clientWidth;
        // Если ширина по какой-то причине равна 0 (контейнер скрыт), предохраняемся от NaN
        if (width === 0) return; 

        // Чистый и безопасный расчет текущего индекса слайда
        currentIndex = Math.round(target.scrollLeft / width);
    }

    $effect(() => {
        autoplayInterval = setInterval(() => {
            if (!autoplayDestroyed && project.items.length > 1) {
                nextSlide();
            }
        }, 5000);

        return () => {
            if (autoplayInterval) clearInterval(autoplayInterval);
        };
    });
</script>

<section id={project.id} class="w-full py-16 border-b border-neutral-900 scroll-mt-20">
<div class="container mx-auto px-4 mb-8">
        <div class="max-w-4xl relative">
            <!-- Заголовок-ссылка -->
            <button 
                onclick={copyProjectLink}
                class="group flex items-center gap-3 text-left focus:outline-none"
                title="Копировать ссылку на проект"
            >
                <h2 class="text-3xl font-bold text-neutral-100 group-hover:text-amber-500 transition-colors">
                    {project.title}
                </h2>
                
                <!-- Иконка якоря, которая появляется при наведении -->
                <span class="text-neutral-600 group-hover:text-amber-500 opacity-0 group-hover:opacity-100 transition-all transform translate-x-[-10px] group-hover:translate-x-0">
                    🔗
                </span>

                <!-- Всплывающая подсказка о копировании -->
                {#if showCopyTooltip}
                    <span class="absolute -top-8 left-0 bg-amber-500 text-neutral-950 text-[10px] font-bold uppercase py-1 px-2 rounded animate-bounce">
                        Ссылка скопирована!
                    </span>
                {/if}
            </button>
            <p class="text-neutral-400 font-light">{project.description}</p>
        </div>
    </div>

    <div class="relative w-full group">
        <div 
            bind:this={sliderElement}
            onscroll={handleScroll}
            class="w-full flex overflow-x-auto snap-x snap-mandatory no-scrollbar scrollbar-none"
        >
            {#each project.items as item}
                <div class="w-full flex-shrink-0 snap-start snap-always flex justify-center px-4">
                    <!-- КОНТЕЙНЕР С ОГРАНИЧЕНИЕМ ШИРИНЫ -->
                    <div class="w-full max-w-5xl mx-auto h-auto">
                        <CarouselItem {item} />
                    </div>
                </div>
            {/each}
        </div>

        {#if project.items.length > 1}
            <button 
                onclick={prevSlide}
                class="hidden md:flex absolute left-4 top-1/2 -translate-y-1/2 w-12 h-12 items-center justify-center bg-black/50 text-white rounded-full opacity-0 group-hover:opacity-100 transition-opacity hover:bg-black/80"
                aria-label="Previous slide"
            >  ← </button>
            <button 
                onclick={() => { stopAutoplayForever(); nextSlide(); }}
                class="hidden md:flex absolute right-4 top-1/2 -translate-y-1/2 w-12 h-12 items-center justify-center bg-black/50 text-white rounded-full opacity-0 group-hover:opacity-100 transition-opacity hover:bg-black/80"
                aria-label="Next slide"
            > → </button>
        {/if}
    </div>

    <div class="container mx-auto px-4 mt-4 flex items-center justify-between">
        <div class="flex gap-1.5">
            {#each project.items as _, i}
                <button onclick={() => scrollToIndex(i)} class="h-1 transition-all rounded-full {i === currentIndex ? 'w-8 bg-amber-500' : 'w-2 bg-neutral-700'}" aria-label="Go to slide {i + 1}"></button>
            {/each}
        </div>
        <div class="text-xs tracking-widest text-neutral-600 font-mono">{currentIndex + 1} / {project.items.length}</div>
    </div>
</section>

<style>
    .scrollbar-none::-webkit-scrollbar {
        display: none;
    }
    @keyframes bounce {
        0%, 100% { transform: translateY(0); }
        50% { transform: translateY(-5px); }
    }
    .animate-bounce {
        animation: bounce 0.5s ease-in-out infinite;
    }
</style>
<script lang="ts">
    import type { Review } from '$lib/types';
    import ReviewBlock from './ReviewBlock.svelte';

    let { reviews }: { reviews: Review[] } = $props();
</script>

<section id="reviews" class="w-full py-24 bg-neutral-950 border-t border-neutral-900 scroll-mt-16">
    <div class="container mx-auto px-4 max-w-6xl">
        <h2 class="text-4xl font-bold text-neutral-100 mb-12">Отзывы</h2>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {#each reviews as review}
                <ReviewBlock {review} />
            {/each}
        </div>
    </div>
</section>

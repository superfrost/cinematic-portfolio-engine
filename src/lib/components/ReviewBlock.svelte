<script lang="ts">
    import type { Review } from '$lib/types';
    import { getEmbedUrl } from '$lib/utils';

    let { review }: { review: Review } = $props();
</script>

<div class="bg-neutral-900/40 p-6 rounded-xl border border-neutral-800 flex flex-col gap-5 hover:border-amber-500/50 transition-colors">
    <div class="aspect-video bg-black rounded-lg overflow-hidden border border-neutral-800 shadow-inner">
        {#if review.video?.provider && review.video?.id}
            <iframe
                class="w-full h-full border-0"
                src={getEmbedUrl(review.video.provider, review.video.id)}
                title="Review video"
                allowfullscreen
            ></iframe>
        {:else if review.image}
            <img src={review.image} alt="Review" class="w-full h-full object-cover" />
        {/if}
    </div>

    <div>
        <div class="flex items-center gap-2 mb-2 text-amber-500">
            {#each Array(5) as _}<span>★</span>{/each}
        </div>
        <h3 class="text-xl font-bold text-neutral-100">{review.author}</h3>
        <p class="text-sm text-neutral-500 mb-3">{review.projectType}</p>
        <p class="text-neutral-400 font-light leading-relaxed italic">«{review.text}»</p>
    </div>
</div>
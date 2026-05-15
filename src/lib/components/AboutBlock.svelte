<script lang="ts">
    import type { Author } from '$lib/types';
    import { resolve } from '$app/paths';
    import { getEmbedUrl } from '$lib/utils';

    let { author }: { author: Author } = $props();
</script>

<div class="w-full flex items-center bg-neutral-900/20 py-20 border-t border-neutral-900 scroll-mt-16">
    <div class="container mx-auto px-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-12 items-start max-w-5xl mx-auto">

            <div class="flex flex-col gap-6 w-full">
                {#if author.media?.photo}
                    <div class="w-full bg-neutral-800 rounded-2xl overflow-hidden shadow-2xl border border-neutral-800 group">
                        <img
                            src={resolve(author.media.photo)}
                            alt={author.name}
                            class="w-full h-auto grayscale contrast-125 filter group-hover:grayscale-0 transition-all duration-700"
                        />
                    </div>
                {/if}

                {#if author.media?.video}
                    <div class="aspect-video bg-neutral-950 rounded-2xl overflow-hidden border border-neutral-800 shadow-lg w-full">
                        <iframe
                            class="w-full h-full border-0"
                            src={getEmbedUrl(author.media.video.provider, author.media.video.id)}
                            title="Author showreel"
                            allowfullscreen
                        ></iframe>
                    </div>
                {/if}

                {#if !author.media?.photo && !author.media?.video}
                    <div class="p-12 border-2 border-dashed border-neutral-800 rounded-2xl text-neutral-600 text-center">
                        Медиа-контент не указан
                    </div>
                {/if}
            </div>

            <div class="flex flex-col justify-center sticky top-24">
                <span class="text-xs uppercase tracking-widest text-amber-500 font-semibold mb-3 block">
                    {author.profession}
                </span>
                <h2 class="text-3xl md:text-5xl font-bold tracking-tight text-neutral-100 mb-6">
                    {author.name}
                </h2>
                <p class="text-neutral-400 font-light text-base md:text-lg leading-relaxed mb-8">
                    {author.bio}
                </p>

                <div class="space-y-4 border-t border-neutral-800 pt-6 font-light">
                    {#if author.contacts.telegram}
                        <div class="flex justify-between py-2 border-b border-neutral-900">
                            <span class="text-neutral-500 text-sm">Telegram:</span>
                            <a href="https://t.me/{author.contacts.telegram.replace('@', '')}" class="hover:text-amber-500 transition-colors">
                                {author.contacts.telegram}
                            </a>
                        </div>
                    {/if}

                    {#if author.contacts.email}
                        <div class="flex justify-between py-2 border-b border-neutral-900">
                            <span class="text-neutral-500 text-sm">Email:</span>
                            <a href="mailto:{author.contacts.email}" class="hover:text-amber-500 transition-colors">
                                {author.contacts.email}
                            </a>
                        </div>
                    {/if}

                    <div class="flex justify-between py-2 border-b border-neutral-900">
                        <span class="text-neutral-500 text-sm">Локация:</span>
                        <span class="text-neutral-300">{author.location}</span>
                    </div>
                </div>
            </div>

        </div>
    </div>
</div>

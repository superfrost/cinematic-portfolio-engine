<script lang="ts">
    let { item } = $props();
</script>

<div class="w-full flex items-center justify-center py-4">
    {#if item.type === 'video'}
        <div class="w-full aspect-video bg-black shadow-2xl rounded-lg overflow-hidden border border-neutral-800">
            <iframe
                class="w-full h-full"
                src={item.provider === 'youtube' 
                    ? `https://www.youtube.com/embed/${item.id}?rel=0` 
                    : item.provider === 'vimeo' 
                    ? `https://player.vimeo.com/video/${item.id}` 
                    : `https://rutube.ru/play/embed/${item.id}`}
                title="Video"
                allowfullscreen
            ></iframe>
        </div>

    {:else if item.type === 'image'}
        <!-- Картинка просто уменьшается, сохраняя пропорции -->
        <img 
            src={item.src} 
            alt="Portfolio item" 
            class="w-full h-auto max-h-[80vh] object-contain rounded-lg shadow-xl border border-neutral-800"
            loading="lazy"
        />

    {:else if item.type === 'text'}
        <div class="w-full h-auto flex items-center justify-center p-6 md:p-12">
            <div class="max-w-3xl w-full p-8 md:p-16 rounded-2xl 
                        bg-gradient-to-br from-amber-500/20 via-neutral-900 to-neutral-950 
                        border border-amber-500/30 shadow-[0_0_50px_-12px_rgba(245,158,11,0.3)]">
                <p class="text-2xl md:text-5xl font-bold tracking-tight text-neutral-100 leading-tight">
                    {item.content}
                </p>
                <div class="mt-8 w-12 h-1 bg-amber-500"></div>
            </div>
        </div>
    {/if}
</div>
<script lang="ts">
    import { page } from '$app/stores';
</script>

<svelte:head>
    <title>404 — Страница не найдена</title>
</svelte:head>

<div class="min-h-screen bg-neutral-950 text-neutral-100 flex items-center justify-center">
    <div class="text-center">
        <h1 class="text-6xl font-bold text-amber-500 mb-4">404</h1>
        <p class="text-neutral-400">Страница не найдена</p>
        <a href="/" class="text-amber-500 hover:underline mt-4 inline-block">На главную</a>
    </div>
</div>

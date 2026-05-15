<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
    <link rel="icon" href={favicon} />
    <title>Илья Пурясев — Видеограф | продюсер</title>
    <meta name="description" content="Портфолио видеографа и продюсера. Коммерческое видеопроизводство, рекламные промо-ролики и музыкальные клипы." />
    <meta property="og:title" content="Илья Пурясев — Видеограф | продюсер" />
    <meta property="og:description" content="Портфолио видеографа и продюсера. Коммерческое видеопроизводство, рекламные промо-ролики и музыкальные клипы." />
    <meta property="og:type" content="website" />
</svelte:head>
{@render children()}

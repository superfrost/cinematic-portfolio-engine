import type { ComponentProps } from 'svelte';
import PortfolioBlock from '$lib/components/PortfolioBlock.svelte';

type Project = ComponentProps<typeof PortfolioBlock>['project'];

export const prerender = true;

export async function load() {
    // 1. Получаем модули кейсов
    const caseModules = import.meta.glob('$lib/content/cases/*.json', { eager: true });
    
    // 2. Получаем пути к файлам и сортируем их по алфавиту (от А до Я)
    // Затем разворачиваем (reverse), чтобы последние даты (2026...) были первыми
    const sortedCases = Object.keys(caseModules)
        .sort()
        .reverse()
        .map((path) => (caseModules[path] as any).default);

    // 3. То же самое для отзывов (если там важен порядок)
    const reviewModules = import.meta.glob('$lib/content/reviews/*.json', { eager: true });
    const reviews = Object.values(reviewModules).map(m => (m as any).default);

    // Данные автора (импортируем напрямую)
    const authorModule = await import('$lib/content/author.json');

    return { 
        cases: sortedCases,
        reviews: reviews,
        author: authorModule.default
    }
}
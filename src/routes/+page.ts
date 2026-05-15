import type { Project, Review, Author } from '$lib/types';

export const prerender = true;

export async function load() {
    const caseModules = import.meta.glob<{ default: Project }>('$lib/content/cases/*.json', { eager: true });

    const sortedCases = Object.keys(caseModules)
        .sort()
        .reverse()
        .map(path => caseModules[path].default);

    const reviewModules = import.meta.glob<{ default: Review }>('$lib/content/reviews/*.json', { eager: true });
    const reviews = Object.values(reviewModules).map(m => m.default);

    const authorModule = await import('$lib/content/author.json');
    const author = authorModule.default as Author;

    return {
        cases: sortedCases,
        reviews,
        author
    }
}
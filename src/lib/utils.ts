export function getEmbedUrl(provider: string, id: string): string {
    switch (provider) {
        case 'youtube':
            return `https://www.youtube.com/embed/${id}?rel=0`;
        case 'vimeo':
            return `https://player.vimeo.com/video/${id}`;
        case 'rutube':
            return `https://rutube.ru/play/embed/${id}`;
        default:
            return '';
    }
}

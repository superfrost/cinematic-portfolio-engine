export interface MediaItem {
    type: 'video' | 'image' | 'text';
    provider?: 'youtube' | 'vimeo' | 'rutube';
    id?: string;
    src?: string;
    content?: string;
}

export interface Project {
    id: string;
    title: string;
    description: string;
    items: MediaItem[];
}

export interface Review {
    author: string;
    projectType: string;
    text: string;
    video?: {
        provider: string;
        id: string;
    };
    image?: string;
}

export interface Author {
    name: string;
    profession: string;
    bio: string;
    location: string;
    contacts: {
        telegram?: string;
        email?: string;
    };
    media: {
        photo?: string;
        video?: {
            provider: string;
            id: string;
        };
    };
}

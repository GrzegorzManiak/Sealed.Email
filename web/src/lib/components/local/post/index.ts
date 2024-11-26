export type LinkData = {
    title: string;
    url: string;
};

export type PeopleData = {
    first_name: string;
    last_name: string;
    links?: Array<[string, string]>
};

export type PostData = {
    id: string;
    title: string;
    body: string;
    link?: string;
    date: Date;
    stringDate?: string;
    cover: string;
    tags: Array<string>;
    people?: Array<PeopleData>;
    links?: Array<LinkData>;
    markdown?: string;
};

export { default as Post } from './post.svelte';
export { default as PostList } from './list.svelte';
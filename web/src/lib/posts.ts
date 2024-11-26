import {People} from "./people";
import type {PostData} from "$local/post";
import * as JsSearch from 'js-search';
import {
    FeaturedArticles,
    AllArticles,
    PostYears
} from "./generated/posts";
export * from "./generated/posts";



const searcher = new JsSearch.Search('id');
searcher.addIndex('title');
searcher.addIndex('body');
searcher.addIndex('tags');
searcher.addIndex(['people', 'first_name']);
searcher.addIndex(['people', 'last_name']);
searcher.addIndex(['date', 'getFullYear']);
searcher.addIndex('stringDate');
searcher.addDocuments(AllArticles);

function GetRandomPosts(count: number, cutoff: number = 20): Array<PostData> {
    const randomPosts: Array<PostData> = [];
    let i = 0;
    while (randomPosts.length < count) {
        const randomIndex = Math.floor(Math.random() * AllArticles.length);
        const randomPost = AllArticles[randomIndex];
        // @ts-ignore
        if (randomPosts.includes(randomPost) && i++ < cutoff) continue;
        // @ts-ignore
        randomPosts.push(randomPost);
    }
    return randomPosts;
}

const RecentPosts = AllArticles.slice(0, 5);

export {
    PostYears,
    AllArticles as Posts,
    searcher,
    FeaturedArticles as FeaturedPosts,
    RecentPosts,
    GetRandomPosts,
};
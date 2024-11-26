<script lang="ts">

    // -- Start, end and items
    import {Post, type PostData} from "$local/post/index";

    export let start: number = 0;
    export let end: number = 0;
    export let data: Array<PostData> = [];

    $: items = data.slice(start, end);
    let reRender = true;
    $: {
        reRender = !reRender;
        console.log("Re-rendering", items.length);
    }
</script>

{#if reRender}
    {#each items as post}
        <Post
                id={post.id}
                title={post.title}
                body={post.body}
                link={post.link}
                date={post.date}
                stringDate={post.stringDate}
                cover={post.cover}
                tags={post.tags}
                people={post.people}
                links={post.links}
                markdown={post.markdown}
                image={true} heading="text-1xl"/>
    {/each}
{:else}
    {#each items as post}
        <Post
                id={post.id}
                title={post.title}
                body={post.body}
                link={post.link}
                date={post.date}
                stringDate={post.stringDate}
                cover={post.cover}
                tags={post.tags}
                people={post.people}
                links={post.links}
                markdown={post.markdown}
                image={true} heading="text-1xl"/>
    {/each}
{/if}
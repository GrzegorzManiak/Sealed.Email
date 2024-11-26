<script lang="ts">
    import type {PostData} from "$local/post/index";
    import { AspectRatio } from "$shadcn/aspect-ratio";
    import {cn} from "$lib/utils";
    import SvelteMarkdown from 'svelte-markdown';

    type Props = PostData & { class?: string };
    type $$props = Props;
    let className: $$props["class"];
    const props = $$props as Props & { class?: string };

    export let isCard: boolean = false;
    export let heading: string = 'text-2xl';
    export let image: boolean = true;
    export let showTags = true;
    export let trunkAt = 200;
    if (trunkAt > 0) props.body = props.body.slice(0, trunkAt);

    const hasMarkdown = !!props.markdown && trunkAt <= 0;
    const text = hasMarkdown ? props.markdown : props.body;
</script>

<style>
    .sr-only {
        position: absolute;
        width: 1px;
        height: 1px;
        padding: 0;
        margin: -1px;
        overflow: hidden;
        clip: rect(0, 0, 0, 0);
        white-space: nowrap;
        border: 0;
    }
</style>

<div class={cn("flex transition duration-200 w-full", className, isCard ? 'gap-1' : 'gap-4' )}>
    <article class="w-full" id={props.id} data-tags="{props.tags.join(', ')}">
        <header class="w-full flex flex-col gap-0 transition duration-200">

            {#if isCard || !hasMarkdown}
                <a href="/post/{props.id}" class="{heading} font-bold text-primary">{props.title}</a>
            {:else}
                <h1 class="{heading} font-bold text-primary">{props.title}</h1>
            {/if}

            <div class='flex gap-2 justify-start align-middle items-center mb-2'>
                <time class="text-sm text-muted-foreground mt-0" dateTime={props.date.toISOString()}>
                    Date: {props.date.toLocaleDateString()}
                </time>

                {#if props.link}
                    <a
                        class="text-sm text-primary"
                        href={props.link}
                        rel="canonical">View original post here.</a>
                {/if}
            </div>
        </header>


        {#if image && !isCard}
            <div class="w-full { !hasMarkdown ? 'cursor-pointer' : ''}"
                 on:click={() => !hasMarkdown ? window.location.href = `/post/${props.id}` : '' }
                 on:keydown={(e) => {if (e.key === 'Enter' && !hasMarkdown) window.location.href = `/post/${props.id}`}}
                 tabindex="0"
                 role="link">

                <figure>
                    <AspectRatio ratio={16 / 9} class="bg-muted rounded">
                        <img src="{props.cover}" alt="{props.title}" class="object-cover w-full h-full rounded" />
                    </AspectRatio>
                    <figcaption class="text-sm text-muted-foreground mt-2">{props.title}</figcaption>
                </figure>
            </div>
        {/if}



        {#if props.people && props.people.length > 0}
            <section aria-label="Mentioned People" class="sr-only">
                <h2 class="text-lg font-bold">People Mentioned</h2>
                <ul class="flex flex-col gap-2">
                    {#each props.people as person}
                        <li>
                            <span class="font-medium">{person.first_name} {person.last_name}</span>
                            {#if person.links && person.links.length > 0}
                                <ul class="flex flex-row gap-3">
                                    {#each person.links as link}
                                        <li>
                                            <a href="{link[0]}" target="_blank" rel="noopener noreferrer" class="text-blue-600 underline">
                                                {link[1]}
                                            </a>
                                        </li>
                                    {/each}
                                </ul>
                            {/if}
                        </li>
                    {/each}
                </ul>
            </section>
        {/if}

        <section class="flex flex-col { isCard ? '' : 'gap-2' } mt-2 w-full">

            {#if hasMarkdown && image}
                <style>
                    .markdown-container img {
                        max-width: 100%;
                        height: auto;
                        border-radius: 0.25rem;
                    }

                    .markdown-container h1,
                    .markdown-container h2,
                    .markdown-container h3,
                    .markdown-container h4 {
                        margin-top: 1.5rem;
                        font-weight: bold;
                        margin-bottom: 0.5rem;
                        padding-bottom: 0.25rem;
                        border-bottom: 1px solid hsl(var(--border) / var(--tw-border-opacity));
                    }

                    .markdown-container h1 {
                        font-size: 2rem;
                    }

                    .markdown-container h2 {
                        font-size: 1.5rem;
                    }

                    .markdown-container h3 {
                        font-size: 1.25rem;
                    }

                    .markdown-container h4 {
                        font-size: 1rem;
                    }

                    .markdown-container p {
                        margin-top: 1rem;
                        margin-bottom: 1rem;
                    }

                    .markdown-container em {
                        font-style: italic;
                    }

                    .markdown-container strong {
                        font-weight: bold;
                    }

                    .markdown-container blockquote {
                        border-left: 0.2rem solid hsl(var(--muted-foreground) / var(--tw-text-opacity));
                        padding-left: 0.8rem;
                        margin-left: 0;
                        color: hsl(var(--muted-foreground) / var(--tw-text-opacity));
                    }

                    .markdown-container code {
                        background-color: hsl(var(--foreground) / 10%);
                        padding: 0.2rem 0.4rem;
                        border-radius: 0.25rem;
                        font-size: 0.8rem;
                        color: hsl(var(--muted-foreground) / var(--tw-text-opacity));
                    }

                    .markdown-container pre {
                        background-color: hsl(var(--foreground) / 10%);
                        padding: 0.8rem;
                        border-radius: 0.25rem;
                        font-size: 0.8rem;
                        color: hsl(var(--muted-foreground) / var(--tw-text-opacity));
                        overflow-x: auto;
                        margin-top: 1rem;
                        margin-bottom: 1rem;
                    }

                    .markdown-container pre code {
                        background-color: transparent;
                        padding: 0;
                        border-radius: 0;
                        font-size: 0.9rem;
                        color: inherit;
                    }

                    .markdown-container a {
                        color: rgb(37 130 200 / var(--tw-text-opacity));
                        text-decoration: none;
                    }

                    .markdown-container a:hover {
                         text-decoration: underline;
                    }

                    .markdown-container ul {
                        margin: 1rem 0;
                        padding-left: 1.5rem;
                        list-style-position: inside;
                        list-style-type: disc; /* Adds dots */
                    }

                    .markdown-container ol {
                        margin: 1rem 0;
                        padding-left: 1.5rem;
                        list-style-position: inside;
                        list-style-type: decimal; /* Adds numbers */
                    }

                    .markdown-container li {
                        margin: 0rem 0;
                        line-height: 1.5;
                    }

                    .markdown-container table {
                        width: 100%;
                        border-collapse: collapse;
                        margin-top: 1rem;
                        margin-bottom: 1rem;
                    }
                </style>
                <div class="w-full prose max-w-none markdown-container">
                    <SvelteMarkdown source={text} />
                </div>
            {:else}
                <a href="/post/{props.id}" class="body-text text-lg text-muted-foreground prose max-w-none truncate">
                    {props.body}
                </a>
            {/if}

            {#if props.tags.length > 0 && showTags}
                <section aria-label="Tags">
                    <ul class="flex flex-row gap-2 w-full { isCard ? '' : 'mt-1' } mb-1 flex-wrap">
                        <li class="text-sm text-muted-foreground py-1 px-0 rounded select-none transition duration-200" aria-label="Tags:">
                            Tags:
                        </li>
                        {#each props.tags as tag}
                            <li class="text-sm text-muted-foreground py-1 px-0 rounded select-none transition duration-200" aria-label={`Tag: ${tag}`}>
                                <a href='/blog?search="{tag}"'>{tag}</a>
                            </li>
                        {/each}
                    </ul>
                </section>
            {/if}

            {#if trunkAt > 0}
                <div class="relative">
                    <div class="absolute inset-0 flex items-center">
                        <span class="w-full border-t"></span>
                    </div>

                    <div class="relative flex justify-center uppercase">
                    <a
                        href="/post/{props.id}"
                        class="bg-background text-muted-foreground px-2 read-more transition duration-200">
                        READ MORE HERE
                    </a>
                    </div>
                </div>
            {/if}
        </section>

    </article>
</div>

<script lang="ts">
    import {Post, type PostData} from "$local/post/index";
    import {cn} from "$lib/utils";
    import {Author, Publisher} from "$lib/people";
    import {onMount} from "svelte";
    import {GetRandomPosts, RecentPosts} from "$lib/posts";

    import Twitter from "lucide-svelte/icons/twitter";
    import LinkedIn from "lucide-svelte/icons/linkedin";
    import Email from "lucide-svelte/icons/mail";
    import Link from "lucide-svelte/icons/link";


    type Props = PostData & { class?: string };
    type $$props = Props;
    let className: $$props["class"];
    const props = $$props as Props & { class?: string };

    const ArticleJSON = JSON.stringify({
        "@context": "https://schema.org",
        "@type": "NewsArticle",
        "headline": `${props.title}`,
        "image": [
            `${props.cover}`
        ],
        "datePublished": `${props.date}`,
        "author": [ Author ],
        "publisher": Publisher
    });
    
    const ImportantPeopleJSON: Array<string> = [];
    for (const person of props.people ?? []) {
        if (person?.links ?? [].length < 1) continue;
        ImportantPeopleJSON.push(JSON.stringify({
            "@context": "https://schema.org",
            "@type": "Person",
            "name": `${person.first_name} ${person.last_name}`,
            "url": `${person.links ?? [['', '']][0][0]}`,
            "sameAs": person.links ?? [].map(link => link[0])
        }));
    }

    const randomPosts = GetRandomPosts(3);

    const tags = props.tags ?? [];
    const formattedTags = tags.map(tag => tag.replace(/ /g, '-').toLowerCase());
    const peopleTags = (props.people ?? []).map(person => `${person.first_name}-${person.last_name}`.toLowerCase());
    const titleTags = props.title.split(' ').map(tag => tag.toLowerCase());
    const allTags = [...formattedTags, ...peopleTags, ...titleTags];
    const prettyTags = allTags.join(', ');

    onMount(() => {
        const script = document.createElement('script');
        script.type = 'application/ld+json';
        script.innerHTML = ArticleJSON;
        document.head.appendChild(script);
        
        for (const person of ImportantPeopleJSON) {
            const script = document.createElement('script');
            script.type = 'application/ld+json';
            script.innerHTML = person;
            document.head.appendChild(script);
        }
    });
</script>

<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{props.title} - Grzegorz Maniak</title>
    <meta name="description" content="{props.title}. {props.body} - {Author.name}" />

    <meta name="keywords" content="{prettyTags}" />
    <meta name="googlebot-news" content="index, follow, max-image-preview:large, max-snippet:-1, max-video-preview:-1">
    <link rel="canonical" href="https://grzegorz.ie/post/{props.id}" />

    <!-- Open Graph Meta Tags for Social Media -->
    <meta property="og:title" content="{props.title} - Grzegorz Maniak" />
    <meta property="og:description" content="A detailed post about {props.title}. Learn more about {props.title} and the people involved like {Author.name}." />
    <meta property="og:type" content="article" />
    <meta property="og:image" content="{props.cover}" />
    <meta property="og:url" content="https://grzegorz.ie/post/{props.id}" />
    <meta property="og:site_name" content="Grzegorz Maniak" />

    <!-- Twitter Card Meta Tags -->
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:title" content="{props.title} - Grzegorz Maniak" />
    <meta name="twitter:description" content="Discover more about {props.title} and the people involved in this article." />
    <meta name="twitter:image" content="{props.cover}" />

    <!-- Canonical URL -->
    <link rel="canonical" href="https://grzegorz.ie/post/{props.id}" />

    <!-- Favicon -->
    <link rel="icon" href="/favicon.ico" type="image/x-icon" />
</head>


<div class={cn(className, 'w-screen flex justify-center items-center mt-[1rem]')}>
    <div class='max-w-screen-lg w-full h-full flex flex-col gap-4 p-4 md:p-6n align-top'>
        <Post
                showTags={false}
                trunkAt={-1}
                {...props}/>

        <div>
            <div class="relative mb-5">
                <div class="absolute inset-0 flex items-center">
                    <span class="w-full border-t"></span>
                </div>

                <div class="relative flex justify-center uppercase">
                    <span class="bg-background text-muted-foreground text-xs px-2">
                        Share
                    </span>
                </div>
            </div>


            <div class="flex gap-4 justify-center">
                <a
                    href={`https://www.twitter.com/intent/tweet?text=Check out this post on Grzegorz Maniak's blog: ${props.title}&url=https://grzegorz.ie/post/${props.id}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    class="text-gray-400"
                    aria-label="Share on Twitter / X"
                >
                    <Twitter size={24} class="stroke-[1.5px]"/>
                    <p class="sr-only">Share on Twitter</p>
                </a>

                <a
                    href={`https://www.linkedin.com/shareArticle?mini=true&title=${props.title}&summary=${props.body}&url=https://grzegorz.ie/post/${props.id}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    class="text-gray-400"
                    aria-label="Share on LinkedIn"
                >
                    <LinkedIn size={24} class="stroke-[1.5px]"/>
                    <p class="sr-only">Share on LinkedIn</p>
                </a>

                <a
                    href={`mailto:?subject=${props.title}&body=Check out this post on Grzegorz Maniak's blog: https://grzegorz.ie/post/${props.id}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    class="text-gray-400"
                    aria-label="Share via Email"
                >
                    <Email size={24} class="stroke-[1.5px]"/>
                    <p class="sr-only">Share via Email</p>
                </a>

                <a
                    on:click={() => navigator.clipboard.writeText(`https://grzegorz.ie/post/${props.id}`)}
                    href="#"
                    class="text-gray-400"
                    aria-label="Copy Link"
                >
                    <Link size={24} class="stroke-[1.5px]"/>
                    <p class="sr-only">Copy Link</p>
                </a>


            </div>

        </div>

        {#if props.tags.length > 0}
            <section aria-label="Tags">
                <div class="relative">
                    <div class="absolute inset-0 flex items-center">
                        <span class="w-full border-t"></span>
                    </div>

                    <div class="relative flex justify-center uppercase">
                <span class="bg-background text-muted-foreground text-xs px-2">
                    tags
                </span>
                    </div>
                </div>

                <ul class="flex flex-row gap-5 w-full mt-2 flex-wrap">
                    {#each props.tags as tag}
                        <li class="text-sm text-muted-foreground bg-background py-1 px-0 rounded select-none" aria-label={`Tag: ${tag}`}>
                            <a href='/blog?search="{tag}"'>{tag}</a>
                        </li>
                    {/each}
                </ul>
            </section>
        {/if}

        {#if props.people ?? 0 > 0}
            <div class="relative">
                <div class="absolute inset-0 flex items-center">
                    <span class="w-full border-t"></span>
                </div>

                <div class="relative flex justify-center uppercase">
                <span class="bg-background text-muted-foreground text-xs px-2">
                    people
                </span>
                </div>
            </div>

            <section aria-label="People Mentioned" class="flex gap-1 justify-between align-middle flex-wrap">
                {#each props.people ?? [] as person}
                    <div class="person-mention flex flex-col gap-0 justify-start align-middle min-w-[4rem]">
                        {#if person.links ?? 0 > 0}
                            <a
                                href={(person.links ?? [['', '']])[0][0]}
                                target="_blank"
                                rel="noopener noreferrer"
                                class="text-primary text-sm"
                                title={`Visit ${person.first_name} ${person.last_name}'s link`}
                            >
                                {person.first_name} {person.last_name}
                            </a>
                        {:else}
                            <span class="text-primary text-sm">
                                {person.first_name} {person.last_name}
                            </span>
                        {/if}

                        {#if person.links ?? 0 > 0}
                            <ul class="links-list flex flex-col gap-0 max-h-[3rem] overflow-y-auto" aria-label={`Links for ${person.first_name} ${person.last_name}`}>
                                {#each person.links ?? [] as link}
                                    <li>
                                        <a
                                                href={link[0]}
                                                target="_blank"
                                                rel="noopener noreferrer"
                                                class="text-blue-600 underline text-xs"
                                                title={`Visit ${person.first_name} ${person.last_name}'s link`}
                                        >
                                            {link[1]}
                                        </a>
                                    </li>
                                {/each}
                            </ul>
                        {/if}
                    </div>
                {/each}
            </section>
        {/if}

        {#if props.links ?? 0 > 0}
            <div class="relative">
                <div class="absolute inset-0 flex items-center">
                    <span class="w-full border-t"></span>
                </div>

                <div class="relative flex justify-center uppercase">
                <span class="bg-background text-muted-foreground text-xs px-2">
                    links
                </span>
                </div>
            </div>

            <div class="flex flex-col gap-1">
                {#each props.links ?? [] as link}
                    <a
                        href={link.url}
                        class="text-primary text-sm"
                        rel="noopener noreferrer"
                    >
                        {link.title}
                        <span class="text-muted-foreground">({link.url})</span>
                    </a>
                {/each}
            </div>
        {/if}

        <div>
            <div class="relative mb-5">
                <div class="absolute inset-0 flex items-center">
                    <span class="w-full border-t"></span>
                </div>

                <div class="relative flex justify-center uppercase">
                    <span class="bg-background text-muted-foreground text-xs px-2">
                        Random Posts
                    </span>
                </div>
            </div>

            <div class="flex flex-col gap-1">
                <div class='
                    w-full
                    h-full
                    flex
                    flex-col
                    justify-between
                    items-center
                    gap-[2rem]
                    '>

                    <style>
                        .post-small:hover li,
                        .post-small:hover header time,
                        .post-small:hover section a,
                        .post-small:hover a {
                            background-color: inherit;
                            color: black;
                        }

                        .post-small:hover .read-more {
                            background-color: rgb(243 244 246 / var(--tw-bg-opacity));
                        }
                    </style>

                    {#each randomPosts as post}
                        <span class='
                            w-full
                            post-small
                            flex-grow border-l pl-5 transition duration-200
                            hover:bg-gray-100 hover:text-black
                        '>
                            <Post {...post} image={false} heading="text-1xl" isCard={true}/>
                        </span>
                    {/each}
                </div>
            </div>
        </div>

    </div>
</div>

<script lang='ts'>
    import { onMount } from 'svelte';
    import {redirectIfLoggedIn, redirectIfLoggedOut} from "$lib/redirect";
    import {cn} from "$lib/utils";

    import { page } from "$app/stores";
    import { Separator } from '@/ui/separator';
    import { Button } from '@/ui/button';

    import { cubicInOut } from "svelte/easing";
    import { crossfade } from "svelte/transition";

    const sidebarNavItems = [
        {
            title: "Account",
            href: "/inbox/settings/account",
        },
        {
            title: "Security",
            href: "/inbox/settings/security",
        },
        {
        	title: "Appearance",
        	href: "/inbox/settings/appearance",
        },
        {
            title: "Notifications",
            href: "/inbox/settings/notifications",
        },
        {
            title: "Preferences",
            href: "/inbox/settings/preferences",
        },
    ];

    const [send, receive] = crossfade({
        duration: 250,
        easing: cubicInOut,
    });

    onMount(() => {
        redirectIfLoggedOut('/authentication/login');
    });
</script>

<div class='w-full h-full flex justify-center align-top'>
    <div class='flex flex-col justify-between items-center w-full max-w-[60rem] flex-1 h-full'>
        <div class="w-full mt-5">
            <div class="w-full" >
                <h2 class="text-2xl font-bold tracking-tight">Settings</h2>
                <p class="text-muted-foreground">
                    Manage your account settings, preferences, and security here.
                </p>
            </div>

            <Separator class="mt-2"/>
        </div>

        <div class='w-full h-full flex justify-between items-center mt-2 gap-2'>

            <nav class='flex flex-col gap-2 max-w-[15rem] w-full h-full'>
                {#each sidebarNavItems as item}
                    {@const isActive = $page.url.pathname === item.href}

                    <Button
                            href={item.href}
                            variant='ghost'
                            class={cn(!isActive && 'hover:underline', 'relative justify-start hover:bg-transparent w-full')}
                            data-sveltekit-noscroll
                    >

                        {#if isActive}
                            <div
                                    class='absolute inset-0 rounded-md bg-muted'
                                     in:send={{ key: 'active-sidebar-tab' }}
                                    out:receive={{ key: 'active-sidebar-tab' }}
                            />
                        {/if}

                        <div class='relative'>
                            {item.title}
                        </div>
                    </Button>
                {/each}
            </nav>

            <div class='flex-1 w-full h-full'>
                <slot />
            </div>
        </div>
    </div>
</div>

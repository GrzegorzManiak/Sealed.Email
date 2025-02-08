<script lang="ts">
    import {Button} from "@/ui/button";
    import type {Writable} from "svelte/store";
    import {cn} from "$lib/utils";

    export let text: string;
    export let icon: any;

    export let buttonID: string;
    export let stateManager: Writable<string>;
	export let hasNotifications: boolean = false;
	export let collapsed: Writable<boolean>;
    export let notificationId: string = buttonID;
    export let className: string | undefined | null = undefined;
    export let href: string | undefined | null = undefined;

    function handleClick() {
        stateManager.set(buttonID);
    }

	let notifications = Math.floor(Math.random() * 125);
	const ninetyNinePlus = hasNotifications && notifications > 99;
    $: isCurrent = $stateManager === buttonID;
</script>

<Button on:click={handleClick} variant="ghost" {href} class={cn(
    "hover:bg-zinc-800 w-full h-auto rounded-none py-[0.6rem] relative", {
        "bg-zinc-900 text-foreground": isCurrent,
        "hover:bg-opacity-90": !isCurrent
    }, className)}>

    {#if isCurrent}
         <span class="absolute inset-y-0 left-0 w-1 bg-blue-500"></span>
    {/if}

    {#if $collapsed}
        <svelte:component this={icon} size="18" class={cn({ "text-muted-foreground": !isCurrent }, "w-full")}/>
    {:else}
        <div class="w-full flex justify-start items-center align-middle gap-3 m-0 ">
            <svelte:component this={icon} size="18" class={cn({ "text-muted-foreground": !isCurrent })}/>

            <p class={cn({ "text-foreground font-medium": !isCurrent })}>{text}</p>

            <div class="flex-grow"/>

            {#if hasNotifications}
                <p class={cn({ "text-muted-foreground": !isCurrent })}>{ninetyNinePlus ? '99+' : notifications}</p>
            {/if}
        </div>
    {/if}
</Button>
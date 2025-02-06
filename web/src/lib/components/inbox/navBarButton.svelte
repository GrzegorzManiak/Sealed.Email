<script lang="ts">
    import {Button} from "@/ui/button";
    import type {Writable} from "svelte/store";
    import {cn} from "$lib/utils";

    export let text: string;
    export let icon: any;

    export let buttonID: string;
    export let stateManager: Writable<string>;
	export let hasNotifications: boolean = false;

    export let className: string | undefined | null = undefined;

    function handleClick() {
        stateManager.set(buttonID);
    }

    $: isCurrent = $stateManager === buttonID;
</script>

<Button on:click={handleClick} variant="ghost" class={cn("w-full px-[0.75rem] py-[0.475rem] rounded-md h-auto", {
    "bg-blue-600 text-foreground hover:bg-blue-700": isCurrent,
    "bg-transparent hover:bg-card bg-opacity-15 hover:bg-opacity-30": !isCurrent
}, className)}>
    <div class="w-full flex justify-start items-center align-middle gap-3 m-0">
        <svelte:component this={icon} size="16" class={cn({ "text-muted-foreground": !isCurrent }, className)}/>
        <p class={cn({ "text-foreground font-medium": !isCurrent })}>{text}</p>

        <div class="flex-grow"/>

        {#if hasNotifications}
            <p class={cn({ "text-muted-foreground": !isCurrent })}>13</p>
        {/if}
    </div>
</Button>
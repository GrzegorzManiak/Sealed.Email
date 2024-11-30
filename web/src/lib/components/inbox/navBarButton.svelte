<script lang="ts">
    import {Button} from "@/ui/button";
    import type {Writable} from "svelte/store";
    import {cn} from "$lib/utils";

    export let text: string;
    export let icon: any;

    export let buttonID: string;
    export let stateManager: Writable<string>;

    export let className: string | undefined | null = undefined;

    function handleClick() {
        stateManager.set(buttonID);
    }

    $: isCurrent = $stateManager === buttonID;
</script>

<Button on:click={handleClick} variant="ghost" class={cn("w-full px-2 py-2 h-auto", {
    "bg-background bg-opacity-30 hover:bg-background hover:bg-opacity-30": isCurrent,
    "bg-transparent hover:bg-card bg-opacity-15 hover:bg-opacity-15": !isCurrent
}, className)}>
    <div class="w-full flex justify-start items-center align-middle gap-2 m-0">
        <svelte:component this={icon} size="16" class={cn({ "text-foreground text-opacity-75": !isCurrent }, className)}/>
        <p class={cn({ "text-foreground text-opacity-75": !isCurrent })}>{text}</p>
    </div>
</Button>
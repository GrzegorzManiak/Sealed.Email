<script lang="ts">
    import * as Collapsible from "$shadcn/collapsible";
	import {ChevronsUp, ChevronsUpDown, ChevronUp, Globe, Search, Settings, SquarePen} from "lucide-svelte";
    import {Button} from "@/ui/button";
    import NavBarHeader from "@/inbox/navBarHeader.svelte";
    import NavBarButton from "@/inbox/navBarButton.svelte";
	import {writable, type Writable} from "svelte/store";
	import {cn} from "$lib/utils";

	export let collapsed: Writable<boolean>;
	export let text: string;
    export let defaultOpen: boolean = false;

	let open: Writable<boolean> = writable(defaultOpen);
</script>

<Collapsible.Root class="w-full" bind:open={$open}>
    <Collapsible.Trigger asChild let:builder>
        {#if $collapsed}
            <Button builders={[builder]} variant="ghost" size="sm" class="flex flex-col w-full items-center p-1 rounded-none border-t gap-1">
                <span class="sr-only">Toggle</span>
                <NavBarHeader text={text} />
                <div class={cn("h-[3px] rounded-md w-full bg-opacity-40 bg-primary transition-colors", {
                    "bg-opacity-15": $open
                })}/>
            </Button>
        {:else}
            <Button builders={[builder]} variant="ghost" size="sm" class="pl-2 flex justify-between w-full items-center rounded-none border-t">
                <NavBarHeader text={text} className="font-bold text-md" />

                <ChevronUp class={cn({
                    "transform rotate-180": $open
                }, "h-4 w-4 text-muted-foreground transition-transform")}/>

                <span class="sr-only">Toggle</span>
            </Button>
        {/if}
    </Collapsible.Trigger>

    <Collapsible.Content>
        <div class="flex flex-col align-middle items-start justify-center">
            <slot/>
        </div>
    </Collapsible.Content>
</Collapsible.Root>
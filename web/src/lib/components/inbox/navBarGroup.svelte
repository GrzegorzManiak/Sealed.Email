<script lang="ts">
    import * as Collapsible from "$shadcn/collapsible";
    import {ChevronsUpDown, Globe, Search, Settings, SquarePen} from "lucide-svelte";
    import {Button} from "@/ui/button";
    import NavBarHeader from "@/inbox/navBarHeader.svelte";
    import NavBarButton from "@/inbox/navBarButton.svelte";

    export let text: string;
    export let defaultOpen: boolean = false;
</script>

<Collapsible.Root class="w-full" open={defaultOpen}>
    <Collapsible.Trigger asChild let:builder>
        <Button builders={[builder]} variant="ghost" size="sm" class="pl-2 flex justify-between w-full items-center">
            <NavBarHeader text={text} />
            <ChevronsUpDown class="h-4 w-4 text-muted-foreground"/>
            <span class="sr-only">Toggle</span>
        </Button>
    </Collapsible.Trigger>

    <Collapsible.Content>
        <div class="flex flex-col px-3">
            <slot/>
        </div>
    </Collapsible.Content>
</Collapsible.Root>
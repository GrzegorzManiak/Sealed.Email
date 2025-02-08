<script lang="ts">
    import {Inbox, Mails, SquarePen, ChevronLeft} from "lucide-svelte";
    import {Label} from "@/ui/label";
    import {Input} from "@/ui/input";
    import {Button} from "@/ui/button";
    import {Checkbox} from "@/ui/checkbox";
	import type {Writable} from "svelte/store";
	import type {InboxManager} from "$lib/inboxManager";

	export let collapsed: Writable<boolean>;
	export let inboxManager: InboxManager;
</script>

<div class="h-[181px] max-h-[181px] overflow-hidden">
    <div class="flex flex-col items-start justify-between align-bottom p-4 gap-2 border-b bg-primary-foreground bg-opacity-40 h-full">

        <div class="flex flex-col gap-1">
            <h1 class="text-2xl font-bold leading-tight p-0">Encrypted Inbox</h1>
            <div class="flex flex-row gap-2">

                <div class="flex flex-row gap-1 items-center">
                    <Mails class="text-muted-foreground" size={15}/>
                    <p class="text-sm text-muted-foreground leading-tight">
                        1,400 Messages,
                    </p>
                </div>

                <div class="flex flex-row gap-1 items-center">
                    <Inbox class="text-muted-foreground" size={15}/>
                    <p class="text-sm text-muted-foreground leading-tight">
                        216 Unread
                    </p>
                </div>

            </div>
        </div>

        <div class="flex flex-row gap-2 w-full pt-5">
            <Input id="search" type="text" placeholder="Search" class="w-full rounded-md bg-primary-foreground bg-opacity-40 text-foreground"/>
            <Button variant="ghost" class="bg-foreground rounded-md ">
                <SquarePen class="h-5 w-5 text-muted"/>
            </Button>
        </div>

        <div class="flex justify-between align-middle items-center gap-2 w-full">
            <div class="flex flex-row gap-1 items-center">
                <Checkbox id="selectAll" class="text-sm text-muted-foreground border-"/>
                <Label for="selectAll" class="sr-only">Select All</Label>
                <p class="text-sm text-muted-foreground leading-tight">
                    Select All
                </p>
            </div>

            <div class="flex flex-row items-center">
                <p class="text-sm text-muted-foreground leading-tight mr-1">
                    Showing 1-50 of 1,400
                </p>

                <Button variant="ghost" class="w-5 h-5 p-0">
                    <ChevronLeft class="text-muted-foreground"/>
                </Button>

                <Button variant="ghost" class="w-5 h-5 p-0">
                    <ChevronLeft class="text-muted-foreground transform rotate-180"/>
                </Button>
            </div>
        </div>
    </div>
</div>
<script lang='ts'>
    import Pin from "lucide-svelte/icons/pin";
    import Trash from "lucide-svelte/icons/trash";
    import Star from "lucide-svelte/icons/star";
    import MailOpen from "lucide-svelte/icons/mail-open";
    import MailClosed from "lucide-svelte/icons/mail";
    import ChevronRight from "lucide-svelte/icons/chevron-right";
    import ChevronDown from "lucide-svelte/icons/chevron-down";
    import { ActiveDeviceType } from "$lib";
    import type {Attachment, ChainEntry} from "./types";
    import {type InboxContext, InboxManager} from "$lib/inboxManager";
    import {cn} from "$lib/utils";
    import {Checkbox} from "$shadcn/checkbox";

    import * as Avatar from "$shadcn/avatar";
    import * as Tooltip from "$shadcn/tooltip";
    import {Button} from "$shadcn/button";

    import {PrettyPrintTime} from "$lib/time";
    import {TrimText} from "$lib/text";
    import {RandomHEXColor} from "$lib/common";
    import {GetIcon} from "$lib/icons";
    import {Sizes} from "./index";

    type $$props = ChainEntry & { class?: string };
    let className: $$props["class"];
    export { className as class };
    const props = $$props as ChainEntry & { class?: string };
    export let manger: InboxManager;
    export let parentManger: InboxManager;
    export let parentID: string;

    // -- Style & Size
    const device = $ActiveDeviceType as keyof typeof Sizes;
    const sizes = Sizes[device];

    // -- Email state
    let mouseOver = false;
    let isSelected = false;
    let selectMode = false;
    let isClicked = false;

    $: highlighted = mouseOver || selectMode || isClicked;
    $: highlightedColor =
        highlighted && (isClicked || isSelected) ? 'bg-[#0f1524]' :
        highlighted ? 'bg-[#171312]' : ''

    // -- Email actions
    let read = props.read;
    let pinned = false;

    function ToggleRead() {
        read = !read;
    }

    function TogglePinned() {
        pinned = !pinned;
    }

    function ToggleSelected() {
        if (selectMode) return;
        parentManger.select(props.id);
        if (!read) ToggleRead();
    }

    // -- Email data
    let name = props.name || props.email;
    let date = PrettyPrintTime(new Date(props.date));
    let body = TrimText(props.body, 50);
    let color = RandomHEXColor();
    let showAttachment = !!props.showcaseAttachment;
    let showcaseAttachment = props.showcaseAttachment as Attachment;
    let totalAttachments = props.totalAttachments || 1;

    // -- Context
    let context: InboxContext = manger.rawContext;
    manger.context.subscribe((value) => {
        context = value;
        selectMode = context.selectMode;
        // isSelected = value.grouped.has(props.id);
        // if (!value.grouped.has(parentID)) return
        // else isSelected = true;
    });

    parentManger.context.subscribe((value) => {
        isClicked = value.selected === parentID || value.selected === props.id;
    });
</script>

<div
        class={cn(
        "bg-background flex items-stretch justify-between gap-2 select-none cursor-default transition-colors duration-200 relative border-b border-border ",
        highlightedColor
    )}
        on:mouseenter={() => (mouseOver = true)}
        on:mouseleave={() => (mouseOver = false)}
        role="button"
        tabindex="0">


    <!-- Read Indicator -->
    <div class="absolute top-0 left-0 w-1 h-full transition-colors duration-200 {read? 'bg-transparent' : 'bg-blue-500'}"></div>

    <!-- Spacer -->
    <div class="py-4 pl-3 w-8"></div>

    <!-- Left section -->
    <div
            on:click={ToggleSelected}
            on:keydown={(e) => e.key === "Enter" && ToggleSelected()}
            role="button"
            tabindex="0"
            class="py-4"
    >
        <!-- Checkbox (Select Mode) -->
        <div class="w-[40px] h-[40px] flex justify-center items-center transition-opacity {highlightedColor} ">
            <Checkbox bind:checked={isSelected} on:click={() => manger.check(props.id)} aria-label="Select email"/>
        </div>
    </div>



    <!-- Mid section -->
    <div
            on:click={ToggleSelected}
            on:keydown={(e) => e.key === "Enter" && ToggleSelected()}
            role="button"
            tabindex="0"
            style="max-width: {sizes.middleSectionWidth}; min-width: {sizes.middleSectionWidth}"
            class="flex flex-col py-4"
    >
        <div class="flex justify-between">
            <div class="flex items-start justify-between w-full">
                <span class="font-bold text-white text-[{sizes.contactFontSize}]">{name}</span>
                <span class="text-sm text-gray-300 text-[{sizes.subjectFontSize}]">{date}</span>
            </div>
        </div>

        <div class="flex items-start justify-between overflow-clip my-[0.2rem]">
            <span class="text-nowrap w-[70%] truncate max-w-full text-gray-400 text-sm text-[{sizes.bodyFontSize}]">
                {body}
            </span>

            <!-- Read icon -->
            <Tooltip.Root>
                <Tooltip.Trigger>
                <span
                        on:click={() => ToggleRead()}
                        on:keydown={e => e.key === "Enter" && ToggleRead()}
                        role="button"
                        tabindex="0"
                        class="cursor-pointer"
                >
                    {#if read} <MailOpen class="text-gray-400 {mouseOver ? 'opacity-100' : 'opacity-0'} transition-opacity duration-200" size="18"/>
                    {:else} <MailClosed class="text-gray-400 {mouseOver ? 'opacity-100' : 'opacity-0'} transition-opacity duration-200" size="18"/> {/if}
                </span>
                </Tooltip.Trigger>

                <Tooltip.Content> <p>{read? "Mark as unread" : "Mark as read"}</p> </Tooltip.Content>
            </Tooltip.Root>
        </div>

        {#if showAttachment}
            <div class="flex items-center justify-start gap-2 mt-2 max-w-full overflow-clip">
                <Button variant="secondary" size="sm" class="text-sm h-7 max-w-[90%]">
                    <img src={GetIcon(showcaseAttachment.type)} alt={showcaseAttachment.type} class="text-gray-400 max-w-3" />
                    <p class="m-0 ml-1 text-nowrap truncate max-w-full">
                        {showcaseAttachment.name}
                    </p>
                </Button>

                {#if totalAttachments > 1}
                    <span class="text-xs text-gray-400">+{totalAttachments - 1}</span>
                {/if}
            </div>
        {/if}
    </div>



    <!-- Right section -->
    <div class="flex flex-col justify-start align-top items-start gap-2 py-4 pr-4">

        <!-- Pin icon -->
        <Tooltip.Root>
            <Tooltip.Trigger>
                    <span
                            on:click={() => TogglePinned()}
                            on:keydown={(e) => e.key === "Enter" && TogglePinned()}
                            role="button"
                            tabindex="0"
                            class="cursor-pointer"
                    >
                        {#if pinned}<Pin class="text-blue-500 fill-blue-500" size="18" />
                        {:else}<Pin class="text-gray-400 hover:text-blue-300 transition-colors duration-200 {mouseOver ? 'opacity-100' : 'opacity-0'} transition-opacity duration-200" size="18"/>{/if}
                    </span>
            </Tooltip.Trigger>

            <Tooltip.Content> <p>{pinned ? "Unpin" : "Pin"}</p> </Tooltip.Content>
        </Tooltip.Root>


        <!-- Trash icon -->
        <Tooltip.Root>
            <Tooltip.Trigger>
                    <span
                            on:click={() => console.log("Email deleted")}
                            on:keydown={e => e.key === "Enter" && console.log("Email deleted")}
                            role="button"
                            tabindex="0"
                            class="cursor-pointer"
                    >
                        <Trash class="text-gray-400 hover:text-red-500 transition-colors {mouseOver ? 'opacity-100' : 'opacity-0'} transition-opacity duration-200" size="18"/>
                    </span>
            </Tooltip.Trigger>

            <Tooltip.Content> <p>Delete</p> </Tooltip.Content>
        </Tooltip.Root>
    </div>
</div>
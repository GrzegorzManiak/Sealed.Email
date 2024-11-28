<script lang='ts'>
    import Pin from "lucide-svelte/icons/pin";
    import Trash from "lucide-svelte/icons/trash";
    import Star from "lucide-svelte/icons/star";
    import MailOpen from "lucide-svelte/icons/mail-open";
    import MailClosed from "lucide-svelte/icons/mail";
    import ChevronRight from "lucide-svelte/icons/chevron-right";
    import ChevronDown from "lucide-svelte/icons/chevron-down";
    import { ActiveDeviceType } from "$lib";
    import type {Attachment, EmailCardData} from "./types";
    import {type InboxContext, InboxManager} from "$lib/inboxManager";
    import {cn} from "$lib/utils";
    import {Checkbox} from "$shadcn/checkbox";

    import * as Avatar from "$shadcn/avatar";
    import * as Tooltip from "$shadcn/tooltip";

    import {PrettyPrintTime} from "$lib/time";
    import {TrimText} from "$lib/text";
    import {RandomHEXColor} from "$lib/common";
    import {GetIcon} from "$lib/icons";
    import {ChainCard, Sizes} from "./index";
    import {Button} from "$shadcn/button";


    type $$props = EmailCardData & { class?: string };
    let className: $$props["class"];
    export { className as class };
    const props = $$props as EmailCardData & { class?: string };
    const manger = props.inbox;
    const hasChain = props?.chain?.length ?? 0 > 0;

    // -- Style & Size
    const device = $ActiveDeviceType as keyof typeof Sizes;
    const sizes = Sizes[device];

    // -- Email state
    let mouseOver = false;
    let isSelected = false;
    let selectMode = false;
    let isClicked = false;
    let chainVisible = false;

    $: highlighted = mouseOver || selectMode || isSelected || isClicked;
    $: highlightedColor =
        highlighted && (isClicked || isSelected) ? 'bg-[#0f1524]' :
        highlighted ? 'bg-[#01040d]' : ''

    // -- Email actions
    const chainInbox = new InboxManager('chain', 'chain');
    let read = props.read;
    let favorite = props.favorite;
    let pinned = props.pinned;

    function ToggleFavorite() {
        favorite = !favorite;
    }

    function TogglePinned() {
        pinned = !pinned;
    }

    function ToggleRead() {
        read = !read;
    }

    function ToggleChain() {
        chainVisible = !chainVisible;
        manger.chain(props.id);
    }

    function ToggleSelected() {
        if (selectMode) return;
        manger.select(props.id);
        if (!read) ToggleRead();
    }

    function CheckBox() {
        manger.check(props.id);
        if (!hasChain) return;
        chainInbox.check(props.id, true);
    }

    // -- Email data
    let name = props.name || props.email;
    let date = PrettyPrintTime(new Date(props.date));
    let subject = TrimText(props.subject, 50);
    let body = TrimText(props.body, 150);
    let color = RandomHEXColor();
    let avatar = props.avatar;
    let showAttachment = !!props.showcaseAttachment;
    let showcaseAttachment = props.showcaseAttachment as Attachment;
    let totalAttachments = props.totalAttachments || 1;

    // -- Context
    let context: InboxContext = manger.rawContext;
    manger.context.subscribe((value) => {
        context = value;
        isClicked = context.selected === props.id;
        selectMode = context.selectMode;
    });
</script>


<div {...$$restProps}
     class={cn(
        className,
        "bg-background flex-col items-stretch justify-between gap-0 select-none cursor-default transition-colors duration-200 relative w-full",
    )}
>

    <div
        class={cn(
        "bg-background flex items-stretch justify-start gap-2 select-none cursor-default transition-colors duration-200 relative",
        highlightedColor
    )}
        on:mouseenter={() => (mouseOver = true)}
        on:mouseleave={() => (mouseOver = false)}
        role="button"
        tabindex="0">


        <!-- Read Indicator -->
        <div class="absolute top-0 left-0 w-1 h-full transition-colors duration-200 {read? 'bg-transparent' : 'bg-blue-500'}"></div>

        <!-- Email chain indicator -->
        <div class="py-4 pl-3 w-8">
            {#if hasChain}
                <div
                    on:click={() => ToggleChain()}
                    on:keydown={(e) => e.key === "Enter" && ToggleChain()}
                    role="button"
                    tabindex="0"
                    class="cursor-pointer h-[40px] flex items-center"
                >
                    {#if chainVisible} <ChevronDown class="text-gray-400 hover:text-blue-300 transition-colors duration-200" size="18"/>
                    {:else} <ChevronRight class="text-gray-400 hover:text-blue-300 transition-colors duration-200" size="18"/> {/if}
                </div>
            {:else}
                <div class="h-[40px]">
                    <ChevronDown class="text-gray-400 opacity-0" size="18"/>
                </div>
            {/if}
        </div>

        <!-- Left section -->
        <div
            on:click={ToggleSelected}
            on:keydown={(e) => e.key === "Enter" && ToggleSelected()}
            role="button"
            tabindex="0"
            class="flex flex-col justify-start h-full items-center gap-2 py-4"
        >
            <!-- Avatar / Checkbox -->
            <Avatar.Root class="grid grid-cols-1 grid-rows-1 relative">

                <!-- Avatar image -->
                <div class="transition-opacity {highlighted ? 'opacity-0' : 'opacity-100'}">
                    <Avatar.Image style="background-color: {color}" class="select-none" src={avatar} alt={name}/>
                    <Avatar.Fallback style="background-color: {color}">{name}</Avatar.Fallback>
                </div>

                <!-- Checkbox (Select Mode) -->
                <div class="absolute bottom-0 right-0 w-full h-full flex justify-center items-center transition-opacity {highlightedColor} {highlighted ? 'opacity-100' : 'opacity-0'}">
                    <Checkbox bind:checked={isSelected} on:click={() => CheckBox()} aria-label="Select email"/>
                </div>

            </Avatar.Root>


            <!-- Favorite icon -->
            {#if !chainVisible}
                <Tooltip.Root>
                    <Tooltip.Trigger>
                        <span
                            on:click={() => ToggleFavorite()} on:keydown={(e) => e.key === "Enter" && ToggleFavorite()}
                            role="button" tabindex="0" class="cursor-pointer">

                            {#if favorite} <Star class="text-yellow-500 fill-yellow-500 hover:text-yellow-700 hover:fill-yellow-700 transition-colors duration-200" size="18"/>
                            {:else} <Star class="text-gray-500 hover:text-yellow-300 transition-colors duration-200" size="18" /> {/if}
                        </span>
                    </Tooltip.Trigger>

                    <Tooltip.Content>
                        <p>{favorite ? "Remove from favorites" : "Add to favorites"}</p>
                    </Tooltip.Content>
                </Tooltip.Root>
            {/if}
        </div>



        <!-- Mid section -->
        <div
            on:click={ToggleSelected}
            on:keydown={(e) => e.key === "Enter" && ToggleSelected()}
            role="button"
            tabindex="0"
            class="flex flex-col py-4 w-full"
        >
            <div class="flex items-center justify-between">
                <div class="flex items-start justify-between w-max">
                    <span class="font-bold text-white text-[{sizes.contactFontSize}]">{name}</span>
                    <span class="text-sm text-gray-300 text-[{sizes.subjectFontSize}]">{date}</span>
                </div>
            </div>

            <div class="flex flex-col items-start overflow-clip">
                <span class="text-nowrap truncate  flex justify-between  {read ? 'text-gray-300' : 'font-bold text-blue-300'}">
                    {subject}
                    {#if chainVisible}
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
                    {/if}
                </span>

                {#if !chainVisible}
                    <span class="text-nowrap truncate w-max text-gray-400 text-sm text-[{sizes.bodyFontSize}]">{body}</span>
                {/if}
            </div>

            {#if showAttachment && !chainVisible}
                <div class="flex items-center justify-start gap-2 mt-2 max-w-full overflow-clip">
                    <Button variant="secondary" size="sm" class="text-sm h-7 max-w-[90%]">
                        <img src={GetIcon(showcaseAttachment.type)} alt={showcaseAttachment.type} class="text-gray-400 max-w-3" />
                        <p class="m-0 ml-1 text-nowrap truncate max-w-full">{showcaseAttachment.name}</p>
                    </Button>

                    {#if totalAttachments > 1}
                        <span class="text-xs text-gray-400">+{totalAttachments - 1}</span>
                    {/if}
                </div>
            {/if}
        </div>



        <!-- Right section -->
        <div class="flex flex-col justify-start align-middle items-center gap-2 py-4 pr-4">
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


            {#if !chainVisible}
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
            {/if}
        </div>
    </div>

    {#if hasChain && chainVisible}
        <div class="flex flex-col justify-start">

            <ChainCard
                manger={chainInbox}
                parentManger={manger}
                parentID={props.id}
                {...props}
                read={read}
                class="w-full"
            />

            {#each props.chain ?? [] as email}
                <ChainCard
                    manger={chainInbox}
                    parentManger={manger}
                    parentID={props.id}
                    {...email}
                    class="w-full"
                />
            {/each}
        </div>
    {/if}
</div>
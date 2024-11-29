<script lang='ts'>
    import {type Attachment, colors, type Email} from "@/inbox/email";
    import {cn} from "$lib/utils";
    import * as Avatar from "$shadcn/avatar";
    import * as Tooltip from "$shadcn/tooltip";

    import ChevronDown from "lucide-svelte/icons/chevron-down";
    import ChevronRight from "lucide-svelte/icons/chevron-right";
    import Trash from "lucide-svelte/icons/trash";
    import Pin from "lucide-svelte/icons/pin";
    import MailClosed from "lucide-svelte/icons/mail";
    import MailOpen from "lucide-svelte/icons/mail-open";

    import Star from "lucide-svelte/icons/star";
    import {PrettyPrintTime} from "$lib/time";
    import {Checkbox} from "@/ui/checkbox";
    import {GetIcon} from "$lib/icons";
    import {Button} from "@/ui/button";

    export let data: Email;

    let chainVisible = false;
    let favorite = false;
    let isGroupSelected = false;
    let isSelected = false;
    let date = PrettyPrintTime(new Date(data.date));

    function ToggleChain(e: MouseEvent | KeyboardEvent) {
        e.stopPropagation();
        chainVisible = !chainVisible;
    }

    function ToggleFavorite(e: MouseEvent | KeyboardEvent) {
        e.stopPropagation();
        favorite = !favorite;
    }

    function ToggleRead(e: MouseEvent | KeyboardEvent) {
        e.stopPropagation();
        data.read = !data.read;
    }

    function GroupSelect(e: MouseEvent | KeyboardEvent) {
        e.stopPropagation();
        isGroupSelected = !isGroupSelected;
        if (!isGroupSelected) isSelected = false;
    }

    function TogglePinned(e: MouseEvent | KeyboardEvent) {
        e.stopPropagation();
        data.pinned = !data.pinned;
    }

    function ToggleSelected(e: MouseEvent | KeyboardEvent) {
        isSelected = !isSelected;
        data.read = true;
    }


    const avatar = "https://api.dicebear.com/9.x/lorelei/svg?seed=noise.email&options[mood][]=happy";
    let heightOffset = 0;
    let isHovered = false;
</script>


<div
    on:mouseenter={() => isHovered = true}
    on:mouseleave={() => isHovered = false}
    on:click={ToggleSelected}
    on:keydown={(e) => e.key === "Enter" && ToggleSelected(e)}
    role="button"
    tabindex="0"
    class={cn("bg-background flex items-stretch justify-between gap-0 select-none cursor-default transition-colors duration-200 w-full flex-grow relative", {
        [colors.selected]: isSelected || isGroupSelected,
        [colors.hovered]: isHovered && !(isSelected || isGroupSelected),
        [colors.normal]: !isHovered && !(isSelected || isGroupSelected)
    })}>

    <!-- Read Indicator & spacer -->
    <div class="w-1"/>
    <div class={cn({ "bg-blue-500": !data.read }, "absolute h-full transition-colors duration-200 w-1")}/>

    <!-- Chain Indicator & Profile picture & favorite-->
    <div class="flex-col pb-2">

        <!-- Chain Indicator & Profile picture -->
        <div class="flex-col w-min">

            <div class="flex items-center justify-end gap-2 p-2 w-min">
                <!-- Chain Indicator -->
                {#if data.chain && data.chain.length > 0}
                    <div on:click={ToggleChain} on:keydown={(e) => e.key === "Enter" && ToggleChain(e)} role="button" tabindex="0" class="cursor-pointer h-[40px] flex items-center">
                        {#if chainVisible} <ChevronDown class="text-gray-400 hover:text-blue-300 transition-colors duration-200" size="18"/>
                        {:else} <ChevronRight class="text-gray-400 hover:text-blue-300 transition-colors duration-200" size="18"/> {/if}
                    </div>

                {:else}
                    <div class="h-[40px]">
                        <ChevronDown class="text-gray-400 opacity-0" size="18"/>
                    </div>
                {/if}

                <!-- Avatar / Checkbox -->
                <Avatar.Root class={cn("transition-colors duration-200 grid grid-cols-1 grid-rows-1 relative w-10 h-10", {
                        [colors.selected]: isSelected || isGroupSelected,
                        [colors.hovered]: isHovered && !(isSelected || isGroupSelected),
                        [colors.normal]: !isHovered && !(isSelected || isGroupSelected)
                    })}>

                    <!-- Avatar -->
                    <div class={cn("transition-opacity", { 'opacity-0': isHovered || (isSelected || isGroupSelected) })}>
                        <Avatar.Image class="select-none" src={avatar} alt={data.fromName}/>
                        <Avatar.Fallback>{data.fromName}</Avatar.Fallback>
                    </div>

                    <!-- Checkbox (Select Mode) -->
                    <div on:click={(e) => GroupSelect(e)} on:keydown={(e) => e.key === "Enter" && GroupSelect(e)} role="button" tabindex="0" class={cn("absolute bottom-0 right-0 w-full h-full flex justify-center items-center transition-opacity", { 'opacity-0': !isHovered && !(isSelected || isGroupSelected) })}>
                        <Checkbox checked={isGroupSelected} aria-label="Select email"/>
                    </div>
                </Avatar.Root>
            </div>

            <!-- Favorite -->
            <div class="flex justify-end pr-2">
                <div class="w-10 flex justify-center">
                    {#if !chainVisible}
                        <Tooltip.Root>
                            <Tooltip.Trigger>
                                <span on:click={(e) => ToggleFavorite(e)} on:keydown={(e) => e.key === "Enter" && ToggleFavorite(e)} role="button" tabindex="0" class="cursor-pointer">
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
            </div>

        </div>
    </div>


    <!-- Email Content -->
    <div class="max-w-full flex-grow relative mr-2">
        <div class="w-full h-0" style="height: {heightOffset}px"/>


        <!-- Showcase Attachment -->
        {#if data.showcaseAttachment && !chainVisible}
            <div class="flex items-center justify-start gap-2 my-1 max-w-full overflow-clip mb-2">
                <Button variant="secondary" size="sm" class="text-sm h-7 max-w-[90%]">
                    <img src={GetIcon(data.showcaseAttachment.type)} alt={data.showcaseAttachment.type} class="text-gray-400 max-w-3" />
                    <p class="m-0 ml-1 text-nowrap truncate max-w-full">{data.showcaseAttachment.filename}</p>
                </Button>

                {#if (data.totalAttachments ?? 0) > 1}
                    <span class="text-xs text-gray-400">+{data.totalAttachments - 1}</span>
                {/if}
            </div>
        {/if}


        <!-- Email Sender & Date -->
        <div class="absolute left-0 top-0 flex-col justify-start align-baseline max-w-full flex-grow w-full" bind:offsetHeight={heightOffset}>
            <div class="flex items-center justify-between mt-2 w-full">

                <!-- Email Sender -->
                <p class="truncate font-bold text-white">{data.fromName} <span class="truncate font-normal text-sm text-gray-300">&lt;{data.fromEmail}&gt;</span></p>

                <!-- Email actions -->
                <div class="flex items-center gap-2">

                    <!-- Read icon -->
                    {#if isHovered}
                        <Tooltip.Root>
                            <Tooltip.Trigger>
                                <span on:click={(e) => ToggleRead(e)} on:keydown={e => e.key === "Enter" && ToggleRead(e)} role="button" tabindex="0" class="cursor-pointer">
                                    {#if data.read} <MailOpen class="text-gray-400 transition-opacity duration-200" size="18"/>
                                    {:else} <MailClosed class="text-gray-400 transition-opacity duration-200" size="18"/> {/if}
                                </span>
                            </Tooltip.Trigger>
                            <Tooltip.Content> <p>{data.read? "Mark as unread" : "Mark as read"}</p> </Tooltip.Content>
                        </Tooltip.Root>
                    {/if}

                    <!-- Pin icon -->
                    {#if isHovered || data.pinned}
                        <div class="flex items-center gap-2">
                            <Tooltip.Root>
                                <Tooltip.Trigger>
                                <span on:click={(e) => TogglePinned(e)} on:keydown={(e) => e.key === "Enter" && TogglePinned(e)} role="button" tabindex="0" class="cursor-pointer">
                                    {#if data.pinned}<Pin class="text-blue-500 fill-blue-500" size="18" />
                                    {:else}<Pin class="text-gray-400 hover:text-blue-300 transition-colors duration-200" size="18"/> {/if}
                                </span>
                                </Tooltip.Trigger>
                                <Tooltip.Content> <p>{data.pinned ? "Unpin" : "Pin"}</p> </Tooltip.Content>
                            </Tooltip.Root>
                        </div>
                    {/if}
                </div>
            </div>


            <!-- Email Subject & body -->
            <div class="flex-col items-center justify-between w-full">

                <div class="flex items-center justify-between w-full gap-2">
                    <!-- Email Subject -->
                    <p class={cn({"text-gray-300": data.read, "font-bold text-blue-300": !data.read}, "truncate")}>{data.subject}</p>

                    <!-- Email Date -->
                    <p class="text-gray-400 text-sm break-keep whitespace-nowrap">{date}</p>
                </div>

                <!-- Email Body -->
                {#if !chainVisible}<p class="text-gray-400 truncate text-sm">{data.body}</p>{/if}
            </div>
        </div>
    </div>

    <!-- Trash Indicator & spacer -->
    <div class="w-8"/>
    <div class={cn({"opacity-100": isHovered, "opacity-0": !isHovered}, "w-8 h-full absolute right-0")}>
        <Tooltip.Root>
            <Tooltip.Trigger class="w-full transition-colors duration-200 hover:bg-red-500 hover:text-red-800 text-gray-400 flex items-center justify-center h-full">
                <span on:click={() => console.log("Email deleted")} on:keydown={e => e.key === "Enter" && console.log("Email deleted")} role="button" tabindex="0" class="cursor-pointer">
                    <Trash class="text-inherit transition-colors duration-200" size="18"/>
                </span>
            </Tooltip.Trigger>
            <Tooltip.Content> <p>Delete</p> </Tooltip.Content>
        </Tooltip.Root>
    </div>
</div>
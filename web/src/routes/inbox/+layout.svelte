<script lang="ts">
    import NavBarButton from "@/inbox/navBarButton.svelte";
    import StorageIndicator from "@/inbox/storageIndicator.svelte";
    import NavBarHeader from "@/inbox/navBarHeader.svelte";
    import NavBarGroup from "@/inbox/navBarGroup.svelte";
    import * as Avatar from "$shadcn/avatar";
    import * as Collapsible from "$shadcn/collapsible";
	import * as Resizable from '$shadcn/resizable';
	import {writable} from "svelte/store";
	import {Button} from "@/ui/button";
	import {RandomHEXColor} from "$lib/common";
	import {cn} from "$lib/utils";
	import {DomainSelector} from "@/inbox/domainSelector/index.js";

    import {
        Settings,
        Inbox,
        Star,
        Contact,
        Pin,
        SendHorizonal,
        Calendar,
        NotebookPen,
        MailWarning,
        Archive,
        Trash,
    } from "lucide-svelte";

    type Folder = {
        id: string;
        name: string;
        color: string;
    }

    const folders: Array<Folder> = [
        {id: "1", name: "Personal", color: "text-red-300"},
        {id: "2", name: "Work", color: "text-blue-300"},
        {id: "3", name: "School", color: "text-green-200"},
    ];

    type Address = {
        id: string;
        name: string;
    }

    const addresses: Array<Address> = [
        {id: "1", name: "Personal@Grzegorz.ie"},
        {id: "2", name: "Test@Grzegorz.ie"}
    ];

    const stateManager = writable<string>();
	const collapsed = writable(false);

    const color = RandomHEXColor();
    const avatar = `https://api.dicebear.com/9.x/lorelei/svg?seed=${color}icon&options[mood][]=happy`;

	const headerHeight = 190;
</script>

<div class="flex flex-row">
    <Resizable.PaneGroup direction="horizontal" class="h-full">

        <Resizable.Pane
            minSize={20}
            maxSize={25}
            defaultSize={25}
            collapsible={true}
            collapsedSize={7}
            onCollapse={() => collapsed.set(true)}
            onExpand={() => collapsed.set(false)}
            class={cn("flex-shrink-0 h-full", {
                "max-w-[4rem] w-[4rem]": $collapsed
            })}>

            <div class="h-screen flex flex-col gap-1 bg-primary-foreground bg-opacity-40">

                <div class="overflow-y-auto flex-grow flex flex-col">

                    <!-- Compose / Search / Settings -->
                    <div class="pt-2 overflow-hidden" style="height: {headerHeight}px; max-height: {headerHeight}px;">
                        <div class="flex flex-col h-full">
                            <div class="px-2">
                                <DomainSelector />
                            </div>
                            <span class="text-muted-foreground text-xs pb-1"></span>
                            <NavBarButton {collapsed} {stateManager} buttonID="settings" icon={Settings} text="Settings" href="/inbox/settings"/>
                            <NavBarButton {collapsed} {stateManager} buttonID="contacts" icon={Contact} text="Contacts" href="/inbox/contacts"/>
                            <NavBarButton {collapsed} {stateManager} buttonID="calendar" icon={Calendar} text="Calendar" href="/inbox/calendar"/>
                        </div>
                    </div>

                    <!-- Mail -->
                    <NavBarGroup {collapsed} text="Mail" defaultOpen={true}>
                        <NavBarButton {collapsed} hasNotifications={true} {stateManager} buttonID="inbox" icon={Inbox} text="Encrypted Inbox"/>
                        <NavBarButton {collapsed} hasNotifications={true} {stateManager} buttonID="starred" icon={Star} text="Starred"/>
                        <NavBarButton {collapsed} hasNotifications={true} {stateManager} buttonID="pinned" icon={Pin} text="Pinned"/>
                        <NavBarButton {collapsed} hasNotifications={true} {stateManager} buttonID="sent" icon={SendHorizonal} text="Sent"/>
                        <NavBarButton {collapsed} hasNotifications={true} {stateManager} buttonID="later" icon={Calendar} text="Scheduled"/>
                        <NavBarButton {collapsed} hasNotifications={true} {stateManager} buttonID="drafts" icon={NotebookPen} text="Drafts"/>
                        <NavBarButton {collapsed} hasNotifications={true} {stateManager} buttonID="spam" icon={MailWarning} text="Spam"/>
                        <NavBarButton {collapsed} hasNotifications={true} {stateManager} buttonID="archive" icon={Archive} text="Archive"/>
                        <NavBarButton {collapsed} hasNotifications={true} {stateManager} buttonID="trash" icon={Trash} text="Trash"/>
                    </NavBarGroup>


                    <!-- Folders -->
                    <NavBarGroup {collapsed} text="Folders">
                        {#each folders as folder}
                            <NavBarButton {collapsed} {stateManager} buttonID={folder.id} icon={Archive} text={folder.name}/>
                        {/each}

                        {#if folders.length === 0}
                            <p class="text-muted-foreground text-sm text-center">No folders</p>
                        {/if}
                    </NavBarGroup>
                </div>

                <div class="border-t py-2">
                    <StorageIndicator />
                    <div class="flex px-1 mt-1 gap-1 justify-end">

                        <Button href="/logout" variant="secondary" class="text-xs bg-background text-muted-foreground bg-opacity-30 p-1 h-auto">
                            Help
                        </Button>

                        <Button href="/logout" variant="secondary" class="text-xs bg-background text-muted-foreground bg-opacity-30 p-1 h-auto">
                            Feedback
                        </Button>

                        <Button href="/logout" variant="secondary" class="text-xs bg-background text-muted-foreground bg-opacity-30 p-1 h-auto">
                            Logout
                        </Button>
                    </div>
                </div>
            </div>
        </Resizable.Pane>



        <Resizable.Handle withHandle />

        <!-- Content -->
        <slot/>

    </Resizable.PaneGroup>
</div>
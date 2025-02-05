<script lang="ts">
    import NavBarButton from "@/inbox/navBarButton.svelte";
    import StorageIndicator from "@/inbox/storageIndicator.svelte";
    import NavBarHeader from "@/inbox/navBarHeader.svelte";
    import NavBarGroup from "@/inbox/navBarGroup.svelte";
    import * as Avatar from "$shadcn/avatar";
    import * as Collapsible from "$shadcn/collapsible";

    import {
        Search,
        Settings,
        SquarePen,
        Inbox,
        Star,
        Pin,
        SendHorizonal,
        Calendar,
        NotebookPen,
        MailWarning,
        Archive,
        ChevronsUpDown,
        Trash,
        AtSign,
        Globe
    } from "lucide-svelte";

    import {writable} from "svelte/store";
    import {Button} from "@/ui/button";
    import {RandomHEXColor} from "$lib/common";
    import {cn} from "$lib/utils";

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

    const color = RandomHEXColor();
    const avatar = `https://api.dicebear.com/9.x/lorelei/svg?seed=${color}icon&options[mood][]=happy`;

</script>

<div class="flex flex-row bun run ">

    <!-- Left Sidebar -->
    <div class="border-r h-screen flex flex-col gap-1 min-w-[15rem] bg-muted">

        <div class="flex items-center p-2 py-4 gap-2 border-b border-background border-opacity-25">
            <img src="/images/logos/white.png" alt="Noise Logo" class="w-10 mr-1" />
            <h1 class="text-md font-bold">Noise Email</h1>
        </div>

        <div class="overflow-y-auto flex-grow flex flex-col gap-2">

            <!-- Compose / Search / Settings -->
            <div class="border-b border-background border-opacity-25 pb-1">
                <NavBarGroup text="Actions" defaultOpen={true}>
                    <NavBarButton {stateManager} buttonID="compose" icon={SquarePen} text="Compose"/>
                    <NavBarButton {stateManager} buttonID="search" icon={Search} text="Search"/>
                    <NavBarButton {stateManager} buttonID="domains" icon={Globe} text="Domains"/>
                    <NavBarButton {stateManager} buttonID="settings" icon={Settings} text="Settings"/>
                </NavBarGroup>
            </div>

            <!-- Mail -->
            <div class="border-b border-background border-opacity-25 pb-1">
                <NavBarGroup text="Mail" defaultOpen={true}>
                    <NavBarButton {stateManager} buttonID="inbox" icon={Inbox} text="EncryptedInbox"/>
                    <NavBarButton {stateManager} buttonID="starred" icon={Star} text="Starred"/>
                    <NavBarButton {stateManager} buttonID="pinned" icon={Pin} text="Pinned"/>
                    <NavBarButton {stateManager} buttonID="sent" icon={SendHorizonal} text="Sent"/>
                    <NavBarButton {stateManager} buttonID="later" icon={Calendar} text="Scheduled"/>
                    <NavBarButton {stateManager} buttonID="drafts" icon={NotebookPen} text="Drafts"/>
                    <NavBarButton {stateManager} buttonID="spam" icon={MailWarning} text="Spam"/>
                    <NavBarButton {stateManager} buttonID="archive" icon={Archive} text="Archive"/>
                    <NavBarButton {stateManager} buttonID="trash" icon={Trash} text="Trash"/>
                </NavBarGroup>
            </div>

            <!-- Addresses -->
            <div class="border-b border-background border-opacity-25 pb-1">
                <NavBarGroup text="Addresses">
                    {#each addresses as address}
                        <NavBarButton {stateManager} buttonID={address.name} icon={AtSign} text={address.name}/>
                    {/each}

                    {#if addresses.length === 0}
                        <p class="text-muted-foreground text-sm text-center">No addresses</p>
                    {/if}
                </NavBarGroup>
            </div>

            <!-- Folders -->
            <div class="border-b border-background border-opacity-25 pb-1">
                <NavBarGroup text="Folders">
                    {#each folders as folder}
                        <NavBarButton {stateManager} buttonID={folder.id} icon={Archive} text={folder.name}/>
                    {/each}

                    {#if folders.length === 0}
                        <p class="text-muted-foreground text-sm text-center">No folders</p>
                    {/if}
                </NavBarGroup>
            </div>
        </div>

        <div class="border-t border-background border-opacity-25 py-2">
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

<!--            <p class="text-muted-foreground text-xs text-center mt-2">Version hash: ND0FA1238</p>-->
        </div>
    </div>

    <!-- Content -->
    <div class="overflow-clip flex-grow">
        <slot></slot>
    </div>
</div>
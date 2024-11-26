<script lang="ts">
    import { Button } from '$shadcn/button';
    import * as AlertDialog from "$shadcn/alert-dialog";
    import type {Dots, DotsSettings} from "$local/dots/index.d";
    import {DefaultSettings, DotsWritable} from "$local/dots/index";
    import {Input} from "$shadcn/input";
    import {Label} from "$shadcn/label";

    type Props = {} & { class?: string };
    type $$props = Props;
    let className: $$props["class"];
    const props = $$props as Props & { class?: string };

    export let showing = false;

    function SetPopup(visible: boolean) {
        console.log("Setting popup to: ", visible);
        showing = visible;
    }

    function Reset() {
        $DotsWritable = DefaultSettings;
        console.log("Resetting settings");
    }

</script>

<AlertDialog.Root bind:open={showing}>
    <AlertDialog.Content>

        <AlertDialog.Header>
            <AlertDialog.Title>
                Dots Playground
            </AlertDialog.Title>
            <AlertDialog.Description>
                Mess with the dots to your heart's content!
            </AlertDialog.Description>
            <br>



            <div class="flex flex-col gap-2.5 max-h-[40vh] overflow-y-auto">
                <div class="relative">
                    <div class="absolute inset-0 flex items-center">
                        <span class="w-full border-t"></span>
                    </div>

                    <div class="relative flex justify-center text-xs uppercase">
                    <span class="bg-background text-muted-foreground px-2">
                        Text
                    </span>
                    </div>
                </div>

                <div class="flex w-full flex-col gap-1.5">
                    <Label for="text">Text</Label>
                    <Input
                            bind:value={$DotsWritable.text}
                            type="email"
                            id="txt"
                            placeholder="GRZEGORZ" />
                </div>

                <div class="flex w-full flex-col gap-1.5">
                    <Label for="text_height">Text Height</Label>
                    <Input
                            bind:value={$DotsWritable.text_height}
                            type="number"
                            id="text_height"
                            placeholder="15" />
                </div>

                <div class="flex w-full flex-col gap-1.5">
                    <Label for="color">Color</Label>
                    <Input
                            bind:value={$DotsWritable.color}
                            type="text"
                            id="color"
                            placeholder="#000000" />
                </div>


                <div class="relative">
                    <div class="absolute inset-0 flex items-center">
                        <span class="w-full border-t"></span>
                    </div>

                    <div class="relative flex justify-center text-xs uppercase">
                    <span class="bg-background text-muted-foreground px-2">
                        Dots
                    </span>
                    </div>
                </div>


                <div class="flex w-full flex-col gap-1.5">
                    <Label for="dot_size">Dot Size</Label>
                    <Input
                            bind:value={$DotsWritable.dot_size}
                            type="number"
                            id="dot_size"
                            placeholder="15" />
                </div>

                <div class="flex w-full flex-col gap-1.5">
                    <Label for="dot_spacing">Dot Spacing</Label>
                    <Input
                            bind:value={$DotsWritable.dot_spacing}
                            type="number"
                            id="dot_spacing"
                            placeholder="15" />
                </div>

                <div class="relative">
                    <div class="absolute inset-0 flex items-center">
                        <span class="w-full border-t"></span>
                    </div>

                    <div class="relative flex justify-center text-xs uppercase">
                    <span class="bg-background text-muted-foreground px-2">
                        Force
                    </span>
                    </div>
                </div>

                <div class="flex w-full flex-col gap-1.5">
                    <Label for="force">Force</Label>
                    <Input
                            bind:value={$DotsWritable.force}
                            type="number"
                            id="force"
                            placeholder="0.1" />
                </div>

                <div class="flex w-full flex-col gap-1.5">
                    <Label for="force_size">Force Size</Label>
                    <Input
                            bind:value={$DotsWritable.force_size}
                            type="number"
                            id="force_size"
                            placeholder="0.1" />
                </div>

                <div class="flex w-full flex-col gap-1.5">
                    <Label for="max_dist">Max Distance</Label>
                    <Input
                            bind:value={$DotsWritable.max_dist}
                            type="number"
                            id="max_dist"
                            placeholder="100" />
                </div>
            </div>


        </AlertDialog.Header>


        <AlertDialog.Footer>


            <Button on:click={() => SetPopup(false)} class="w-full">
                Close
            </Button>


        </AlertDialog.Footer>
    </AlertDialog.Content>
</AlertDialog.Root>
<script lang="ts">

    import {cn} from "$lib/utils";

    type Props = {} & { class?: string };
    type $$props = Props;
    let className: $$props["class"];
    const props = $$props as Props & { class?: string };

</script>

<div class={cn("w-screen", className)}>
    <div class="flex flex-col gap-4 p-4 md:p-6 w-full max-w-screen-lg mx-auto border-b min-h-[3rem]">

        <div class="w-full h-full flex justify-between items-center align-middle flex-wrap sm:flex-nowrap">
            <div class="flex gap-4 items-center">
                <a
                        href="/"
                        class="md:text-2xl font-bold sm:text-lg"
                >Grzegorz Maniak</a>
            </div>

            <div class="flex items-center sm:gap-2 md:gap-4 gap-2">
                <a href="/blog" class="sm:text-sm md:text-lg text-sm font-medium">Blog</a>
                <a href="/projects" class="sm:text-sm md:text-lg text-sm font-medium">Projects</a>
                <a href="/contact" class="sm:text-sm md:text-lg text-sm font-medium">Contact</a>
                <a href="/about" class="sm:text-sm md:text-lg text-sm font-medium">About</a>
                <a href="/cv" class="sm:text-sm md:text-lg text-sm font-medium">CV</a>
            </div>
        </div>
    </div>
</div>
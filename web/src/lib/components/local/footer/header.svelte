<script lang="ts">

    import {cn} from "$lib/utils";

    type Props = {} & { class?: string };
    type $$props = Props;
    let className: $$props["class"];
    const props = $$props as Props & { class?: string };

</script>

<div class={cn("w-screen mt-[10rem] relative", className)}>

    <div class="
        flex flex-col gap-4 p-4 md:p-6 w-full max-w-screen-lg mx-auto border-t h-[10rem]
    ">

        <div class="flex flex-col gap-5">
            <div class="flex flex-col items-center">
                <div class="text-sm flex gap-3">
                    <a href="mailto:contact@grzegorz.ie" class="text-primary">
                        web@grzegorz.ie
                    </a>

                    <span class="text-muted-foreground"> | </span>

                    <a href="https://www.linkedin.com/in/grzegorz-maniak/" class="text-primary">
                        LinkedIn
                    </a>

                    <span class="text-muted-foreground"> | </span>

                    <a href="https://www.github.com/grzegorzmaniak" class="text-primary">
                        GitHub
                    </a>
                </div>
            </div>
        </div>

        <div class="
            w-full flex justify-between items-center
            flex-wrap
            sm:flex-nowrap
            gap-5
            sm:gap-0
        ">
            <div class="flex flex-col gap-0 items-center sm:items-start w-full">
                <h1 class="text-2xl font-bold">Grzegorz Maniak</h1>
                <p>Software Engineer 🇮🇪 🇵🇱</p>
            </div>

            <div class="flex gap-4 justify-center sm:items-end w-full sm:w-min">
                <a href="/blog" class="text-lg font-medium">Blog</a>
                <a href="/projects" class="text-lg font-medium">Projects</a>
                <a href="/contact" class="text-lg font-medium">Contact</a>
                <a href="/about" class="text-lg font-medium">About</a>
                <a href="/cv" class="text-lg font-medium">CV</a>
            </div>
        </div>

        <div class="w-full">
            <p class="text-sm text-muted-foreground text-center w-full mb-2">© 2024 Grzegorz Maniak</p>
        </div>
    </div>
</div>
<script lang="ts">
    import { onMount } from 'svelte';
    import {redirectIfLoggedIn} from "$lib/redirect";
    import {Button} from "@/ui/button";

	function goBack() {
		history.back();
	}

    onMount(() => {
        redirectIfLoggedIn('/inbox');
    });
</script>

<div
        class="container relative hidden h-screen flex-col items-center justify-center md:grid lg:max-w-none lg:grid-cols-2 lg:px-0"
>

    <div class="relative hidden h-full flex-col bg-muted p-10 text-white dark:border-r lg:flex">


        <!-- Gradient that blocks out the bottom of the splash -->
        <div
                class="absolute inset-0 bg-cover z-10"
                style="background-image: linear-gradient(0deg, rgba(0,0,0,0.75) 0%, rgba(0,0,0,0.45) 15%, rgba(0,0,0,0) 100%);"
        />

        <!-- Splash Image -->
        <div>
            <div class="absolute z-50 flex flex-col">
                <img src="/images/logos/white.png" alt="Noise Logo" class="w-32 mx-auto mr-1" />
            </div>

            <div
                    class="absolute inset-1 bg-cover z-0 top-0 left-0 w-full h-full"
                    style="background-image: url('/images/gradient.png')"
            />
        </div>

        <!-- Quote -->
        <div class="relative z-20 mt-auto">

            <blockquote class="space-y-2">
                <p class="text-lg">
                    True client-to-client, end-to-end encrypted communication, simple and secure.
                </p>
                <footer class="text-sm">Read the blog: <a href="/blog" class="underline hover:text-primary">noise.email/blog</a></footer>
            </blockquote>

        </div>
    </div>

    <div class='absolute top-0 right-0'>
        <div class="p-5">
            <Button variant="ghost" on:click={goBack}>
                Go Back
            </Button>
        </div>
    </div>

    <div class="lg:p-8">
        <slot/>
    </div>
</div>
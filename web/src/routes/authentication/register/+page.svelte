<script lang="ts">
    import { goto } from '$app/navigation';
    import { get } from 'svelte/store';
    import { onMount } from 'svelte';

    import { Button } from '$shadcn/button';
    import { Input } from '$shadcn/input';
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
        <div
                class="absolute inset-0 bg-cover z-0"
                style="background-image: url('/images/splashColor.png')"
        />

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

    <div class="lg:p-8">
        <div class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">

            <!-- Prompt text -->
            <div class="flex flex-col space-y-2 text-center">
                <h1 class="text-2xl font-semibold tracking-tight">Create an account</h1>
                <p class="text-sm text-muted-foreground">
                    Already have an account? <a href="/authentication/login" class="underline hover:text-primary">Click here to Log in.</a>
                </p>
            </div>

            <!-- Form -->

            <!-- Terms -->
            <p class="px-8 text-center text-sm text-muted-foreground">
                By clicking submit, you agree to our
                <a href="/legal/terms" class="underline underline-offset-4 hover:text-primary">
                    Terms of Service
                </a>
                and
                <a href="/legal/privacy" class="underline underline-offset-4 hover:text-primary">
                    Privacy Policy
                </a>
                .
            </p>
        </div>
    </div>
</div>
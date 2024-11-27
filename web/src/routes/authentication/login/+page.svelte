<script lang="ts">
    import { onMount } from 'svelte';
    import { Login } from "$local/login";
</script>

<div class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">

    <!-- Prompt text -->
    <div class="flex flex-col space-y-2 text-center">
        <h1 class="text-2xl font-semibold tracking-tight">Log into your account</h1>
        <p class="text-sm text-muted-foreground">
            Dont have an account? <a href="/authentication/register" class="underline hover:text-primary">Click here to Register.</a>
        </p>
    </div>

    <!-- Form -->
    <Login />

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
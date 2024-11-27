<script lang="ts">
    import { onMount } from 'svelte';
    import { redirectIfLoggedIn } from "$lib/redirect";
    import {Button} from "@/ui/button";
</script>

<div class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">

    <!-- Prompt text -->
    <div class="flex flex-col space-y-2 text-center">
        <h1 class="text-2xl font-semibold tracking-tight">What would you like to do?</h1>
        <p class="text-sm text-muted-foreground">
            Lost your password? <a href="/authentication/reset-password" class="underline hover:text-primary">Click here to reset it.</a>
        </p>
    </div>

    <div class="mx-auto flex w-full justify-center gap-2 sm:w-[350px]">
        <Button class="w-32 h-24" variant="outline" href="/authentication/register"
        >Register</Button>

        <Button class="w-32 h-24" variant="outline" href="/authentication/login"
        >Login</Button>
    </div>

    <!-- Terms -->
    <p class="px-8 text-center text-sm text-muted-foreground">
        By continuing, you agree to our
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
<script lang="ts">
    import {cn, Sleep} from "$lib/utils";
    import {onMount} from "svelte";
    import {throwToast} from "$lib/toasts";

    import LoaderIcon from 'lucide-svelte/icons/loader-circle';

    import { Button } from '$shadcn/button';
    import { Input } from '$shadcn/input';
    import { Label } from '$shadcn/label';
    import { Checkbox } from '$shadcn/checkbox';
    import * as API from "$api/lib";
    import {GenericError} from "$api/lib/errors";

    let className: string | undefined | null = undefined;
    export { className as class };
    let isLoading = false;

    // -- Account Details --
    let username: string = '';
    let usernameError: boolean = false;

    let password: string = '';
    let passwordError: boolean = false;

    function validateDetails() {
        let errorMessages = [];

        usernameError = username.length < 3;
        if (usernameError) errorMessages.push('Username must be at least 3 characters long');

        passwordError = password.length < 8;
        if (passwordError) errorMessages.push('Password must be at least 8 characters long');

        const isError = usernameError || passwordError;
        if (isError) throwToast('You have the following errors', errorMessages.join('; '));
        return isError;
    }

    async function submitDetails() {
        const timeStart = performance.now();
        isLoading = true;
        if (validateDetails()) {
            isLoading = false;
            return;
        }

        async function finish() {
            const timeEnd = performance.now();
            const timeTaken = timeEnd - timeStart;
            await Sleep(1 - timeTaken);
            isLoading = false;
        }

        try {
            const result = await API.Login.Login(username, password);
            const session = new API.Session(result);
            await session.DecryptKeys();
        }

        catch (unknownError) {
            const error = GenericError.from_unknown(unknownError);
            throwToast(error.title, error.message);
            await finish();
            return;
        }

        // -- note: We are not calling finish here because we are redirecting to the login page
        throwToast('Logged in', 'You have been logged in successfully!');
    }

    onMount(() => {
        isLoading = false;
        username = '';
        password = '';

        usernameError = false;
        passwordError = false;
    });
</script>

<div class={cn('grid gap-6', className)} {...$$restProps}>

    <!-- Main form -->
    <form on:submit|preventDefault={submitDetails}>
        <div class='grid gap-2'>

            <!-- Break -->
            <div class='relative'>
                <div class='absolute inset-0 flex items-center'>
                    <span class='w-full border-t' />
                </div>
                <div class='relative flex justify-center text-xs uppercase'>
                    <span class='bg-background px-2 text-muted-foreground'>details</span>
                </div>
            </div>

            <!-- username -->
            <div class='grid gap-1'>
                <Label class='sr-only' for='username'>Username</Label>
                <Input
                        class={cn({ "border-red-500 focus:ring-red-500 focus:border-red-500 bg-red-950 placeholder:text-red-100": usernameError }, "transition-colors")}
                        on:focus={() => usernameError = false}
                        id='username'
                        placeholder='Username'
                        type='text'
                        autocapitalize='none'
                        autocomplete='username'
                        autocorrect='off'
                        disabled={isLoading}
                        bind:value={username}
                />
            </div>

            <!-- Password -->
            <div class='grid gap-1'>
                <Label class='sr-only' for='password'>Password</Label>
                <Input
                        class={cn({ "border-red-500 focus:ring-red-500 focus:border-red-500 bg-red-950 placeholder:text-red-100": passwordError }, "transition-colors")}
                        on:focus={() => passwordError = false}
                        id='password'
                        placeholder='Password'
                        type='password'
                        autocapitalize='none'
                        autocomplete='new-password'
                        autocorrect='off'
                        disabled={isLoading}
                        bind:value={password}
                />
            </div>


            <!-- Submit button -->
            <Button type='submit' disabled={isLoading} class='mt-2'>
                {#if isLoading}
                    <LoaderIcon class='w-5 h-5 animate-spin' />
                {:else}
                    Submit
                {/if}
            </Button>

            <slot/>
        </div>
    </form>

</div>
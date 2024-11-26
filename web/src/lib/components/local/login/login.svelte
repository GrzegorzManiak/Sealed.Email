<script lang="ts">
    import {cn, Sleep} from "$lib/utils";
    import {onMount} from "svelte";
    import {throwToast} from "$lib/toasts";

    import LoaderIcon from 'lucide-svelte/icons/loader-circle';

    import { Button } from '$shadcn/button';
    import { Input } from '$shadcn/input';
    import { Label } from '$shadcn/label';
    import { Checkbox } from '$shadcn/checkbox';

    let className: string | undefined | null = undefined;
    export { className as class };
    let isLoading = false;

    // -- Account Details --
    let nickname: string = '';
    let nicknameError: boolean = false;

    let password: string = '';
    let passwordError: boolean = false;

    function validateDetails() {
        let errorMessages = [];

        nicknameError = nickname.length < 3;
        if (nicknameError) errorMessages.push('Username must be at least 3 characters long');

        passwordError = password.length < 8;
        if (passwordError) errorMessages.push('Password must be at least 8 characters long');

        const isError = nicknameError || passwordError;
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

        console.log({
            nickname,
            password,
        });

        const timeEnd = performance.now();
        const timeTaken = timeEnd - timeStart;
        await Sleep(1500 - timeTaken);
        isLoading = false;
    }

    onMount(() => {
        isLoading = false;
        nickname = '';
        password = '';

        nicknameError = false;
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
                        class={cn({ "border-red-500 focus:ring-red-500 focus:border-red-500 bg-red-950 placeholder:text-red-100": nicknameError }, "transition-colors")}
                        on:focus={() => nicknameError = false}
                        id='username'
                        placeholder='Username'
                        type='text'
                        autocapitalize='none'
                        autocomplete='username'
                        autocorrect='off'
                        disabled={isLoading}
                        bind:value={nickname}
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
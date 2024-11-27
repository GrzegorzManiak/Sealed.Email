<script lang="ts">
    import {cn, Sleep} from "$lib/utils";
    import {onMount} from "svelte";
    import {throwToast} from "$lib/toasts";

    import LoaderIcon from 'lucide-svelte/icons/loader-circle';

    import { Button } from '$shadcn/button';
    import { Input } from '$shadcn/input';
    import { Label } from '$shadcn/label';
    import { Checkbox } from '$shadcn/checkbox';

    import * as API from '$api/lib';

    let className: string | undefined | null = undefined;
    export { className as class };
    let isLoading = false;

    // -- Account Details --
    let username: string = '';
    let usernameError: boolean = false;

    let password: string = '';
    let passwordError: boolean = false;

    let confirmPassword: string = '';
    let confirmPasswordError: boolean = false;

    let acceptedTerms: boolean = false;
    let acceptedTermsError: boolean = false;

    let recoveryEmail: string = '';
    let recoveryEmailError: boolean = false;

    function validateDetails() {
        let errorMessages = [];

        usernameError = username.length < 3;
        if (usernameError) errorMessages.push('Username must be at least 3 characters long');

        passwordError = password.length < 8;
        if (passwordError) errorMessages.push('Password must be at least 8 characters long');

        confirmPasswordError = (password !== confirmPassword) || (confirmPassword.length < 8);
        if (password !== confirmPassword) errorMessages.push('Passwords do not match');

        acceptedTermsError = !acceptedTerms;
        if (acceptedTermsError) errorMessages.push('You must accept the terms and conditions');

        const isError = usernameError || passwordError || acceptedTermsError || confirmPasswordError;
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
            await Sleep(1500 - timeTaken);
            isLoading = false;
        }

        try {
            const result = await API.Register.RegisterUser(username, password);
        }

        catch (error) {
            API.GenericError.fromUnknown(error).toast();
            await finish();
            return;
        }

        // -- note: We are not calling finish here because we are redirecting to the login page
        throwToast('Registered', 'You have been registered successfully!');
    }

    onMount(() => {
        isLoading = false;
        username = '';
        password = '';
        confirmPassword = '';
        acceptedTerms = false;
        recoveryEmail = '';

        usernameError = false;
        passwordError = false;
        confirmPasswordError = false;
        acceptedTermsError = false;
        recoveryEmailError = false;
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
                    <span class='bg-background px-2 text-muted-foreground'>Personal</span>
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

            <!-- Break -->
            <div class='relative'>
                <div class='absolute inset-0 flex items-center'>
                    <span class='w-full border-t' />
                </div>
                <div class='relative flex justify-center text-xs uppercase'>
                    <span class='bg-background px-2 text-muted-foreground'>Security</span>
                </div>
            </div>

            <!-- Email -->
            <div class='grid gap-1'>
                <Label class="sr-only" for="recoveryEmail">Recovery Email</Label>
                <Input
                        class={cn({ "border-red-500 focus:ring-red-500 focus:border-red-500 bg-red-950 placeholder:text-red-100": recoveryEmailError }, "transition-colors")}
                        on:focus={() => recoveryEmailError = false}
                        id='recoveryEmail'
                        placeholder='Recovery Email (Optional)'
                        type='email'
                        autocapitalize='none'
                        autocomplete='email'
                        autocorrect='off'
                        disabled={isLoading}
                        bind:value={recoveryEmail}
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

            <!-- Confirm Password -->
            <div class='grid gap-1'>
                <Label class='sr-only' for='confirmPassword'>Confirm Password</Label>
                <Input
                        class={cn({ "border-red-500 focus:ring-red-500 focus:border-red-500 bg-red-950 placeholder:text-red-100": confirmPasswordError }, "transition-colors")}
                        on:focus={() => confirmPasswordError = false}
                        id='confirmPassword'
                        placeholder='Confirm Password'
                        type='password'
                        autocapitalize='none'
                        autocomplete='new-password'
                        autocorrect='off'
                        disabled={isLoading}
                        bind:value={confirmPassword}
                />
            </div>

            <!-- Break -->
            <div class='relative mt-2'>
                <div class='absolute inset-0 flex items-center'>
                    <span class='w-full border-t' />
                </div>
                <div class='relative flex justify-center text-xs uppercase'>
                    <span class='bg-background px-2 text-muted-foreground'>T&C</span>
                </div>
            </div>

            <div class="flex items-center space-x-2">
                <Checkbox
                        id="terms"
                        bind:checked={acceptedTerms}
                        on:click={() => acceptedTermsError = false}
                        class={cn(
                            {"border-input": !acceptedTermsError},
                            { "border-red-500 focus:ring-red-500 focus:border-red-500 bg-red-950 placeholder:text-red-100": acceptedTermsError },
                            "transition-colors")}
                        disabled={isLoading} />

                <Label for="terms" class="text-sm text-muted-foreground">I accept the <a href="#" class="text-primary">terms and conditions</a></Label>
            </div>

            <!-- Submit button -->
            <Button type='submit' disabled={isLoading} class='mt-2'>
                {#if isLoading}
                    <LoaderIcon class='w-5 h-5 animate-spin' />
                {:else}
                    Register
                {/if}
            </Button>

            <slot/>
        </div>
    </form>

</div>
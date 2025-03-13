<script lang="ts">
    import {cn, Sleep} from "$lib/utils";
    import {onMount} from "svelte";
    import {throwToast} from "$lib/toasts";
    import {user} from "$lib/stores";

    import LoaderIcon from 'lucide-svelte/icons/loader-circle';

    import { Button } from '$shadcn/button';
    import { Input } from '$shadcn/input';
    import { Label } from '$shadcn/label';
    import * as API from "$api/lib";
    import {redirect, redirectIfLoggedIn} from "$lib/redirect";
	import EyeClosed from 'lucide-svelte/icons/eye-closed';
	import EyeOpen from 'lucide-svelte/icons/eye';
    
    let className: string | undefined | null = undefined;
    export { className as class };
    let isLoading = false;

	let errorBoxTitle: string = '';
	let errorBoxMessage: string = '';
	let errorBoxVisible: boolean = false;

	function setErrorBox(title: string, message: string) {
		errorBoxTitle = title;
		errorBoxMessage = message;
		errorBoxVisible = true;
	}

	function closeErrorBox() {
		errorBoxVisible = false;
	}

	// -- Account Details --
    let username: string = '';
    let usernameError: boolean = false;

    let password: string = '';
    let passwordError: boolean = false;

    function validateDetails() {
        let errorMessages = [];

        usernameError = username.length < 3;
        if (usernameError) errorMessages.push('Username must be at least 3 characters long');

        passwordError = password.length < 3;
        if (passwordError) errorMessages.push('Password must be at least 3 characters long');

        const isError = usernameError || passwordError;
        if (isError) setErrorBox('You have the following errors', errorMessages.join('; '));
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
			closeErrorBox();
			API.Session.ClearSessionCookie();
            const result = await API.Login.Login(username, password);
            const session = new API.Session(result);
            await session.DecryptKeys();
            user.set({ isLoggedIn: true, session, username });
            return await redirect('/inbox', {
                delay: 1500,
                title: 'You have been logged in successfully!',
                message: 'You will be redirected to your inbox shortly.'
            });
        }

        catch (error) {
			console.log(error);
			const parsedError =  API.GenericError.fromUnknown(error);
			setErrorBox(parsedError.title, parsedError.message);
            return await finish();
        }
    }

    onMount(() => {
		errorBoxTitle = '';
		errorBoxMessage = '';
		errorBoxVisible = false;

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

            <!-- Error Box -->
            {#if errorBoxVisible}
                <div class='bg-red-950 text-red-300 p-2 rounded-md transition-transform'>
                    <div class='flex items center justify-between '>
                        <h2 class='text-sm font-semibold'>{errorBoxTitle}</h2>
                        <button class='text-red-300' on:click={closeErrorBox} type="button">
                            {#if errorBoxVisible}
                                <EyeOpen class='w-4 h-4' />
                            {:else}
                                <EyeClosed class='w-4 h-4' />
                            {/if}
                        </button>
                    </div>
                    <p class='text-sm'>{
                        errorBoxMessage
                    }</p>
                </div>
            {/if}

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

            <!-- Prompt text -->
            <div class='flex flex-col space-y-2 text-center'>
                <p class='text-sm text-muted-foreground'>
                    Lost your password? <a href='/authentication/reset-password' class='underline hover:text-primary'>Click here to reset it.</a>
                </p>
            </div>
        </div>
    </form>

</div>
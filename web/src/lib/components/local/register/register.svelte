<script lang="ts">
    import {cn} from "$lib/utils";
    import {Icon} from "lucide-svelte";

    import { Button } from '$shadcn/button';
    import { Input } from '$shadcn/input';
    import { Label } from '$shadcn/label';
    import { Checkbox } from '$shadcn/checkbox';

    let className: string | undefined | null = undefined;
    export { className as class };
    let isLoading = false;

    // -- Account Details --
    let nickname: string = '';
    let password: string = '';
    let confirm_password: string = '';
    let accepted_terms: boolean = false;

    async function submit_details() {
        isLoading = true;

        isLoading = false;
    }
</script>

<div class={cn('grid gap-6', className)} {...$$restProps}>

    <!-- Main form -->
    <form on:submit|preventDefault={submit_details}>
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

<!--            &lt;!&ndash; First / Seccond name &ndash;&gt;-->
<!--            <div class='flex gap-2'>-->
<!--                <div class='grid gap-1'>-->
<!--                    <Label class='sr-only' for='first_name'>First Name</Label>-->
<!--                    <Input-->
<!--                            id='first_name'-->
<!--                            placeholder='First Name'-->
<!--                            type='text'-->
<!--                            autocapitalize='words'-->
<!--                            autocomplete='given-name'-->
<!--                            autocorrect='off'-->
<!--                            disabled={isLoading}-->
<!--                            bind:value={first_name}-->
<!--                    />-->
<!--                </div>-->
<!--                <div class='grid gap-1'>-->
<!--                    <Label class='sr-only' for='last_name'>Last Name</Label>-->
<!--                    <Input-->
<!--                            id='last_name'-->
<!--                            placeholder='Last Name'-->
<!--                            type='text'-->
<!--                            autocapitalize='words'-->
<!--                            autocomplete='family-name'-->
<!--                            autocorrect='off'-->
<!--                            disabled={isLoading}-->
<!--                            bind:value={last_name}-->
<!--                    />-->
<!--                </div>-->
<!--            </div>-->

            <!-- username -->
            <div class='grid gap-1'>
                <Label class='sr-only' for='username'>Username</Label>
                <Input
                        id='username'
                        placeholder='username'
                        type='text'
                        autocapitalize='none'
                        autocomplete='username'
                        autocorrect='off'
                        disabled={isLoading}
                        bind:value={nickname}
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
<!--            <div class='grid gap-1'>-->
<!--                <Label class='sr-only' for='email'>Email</Label>-->
<!--                <Input-->
<!--                        id='email'-->
<!--                        placeholder='Name@example.com'-->
<!--                        type='email'-->
<!--                        autocapitalize='none'-->
<!--                        autocomplete='email'-->
<!--                        autocorrect='off'-->
<!--                        disabled={isLoading}-->
<!--                        bind:value={email}-->
<!--                />-->
<!--            </div>-->

            <!-- Password -->
            <div class='grid gap-1'>
                <Label class='sr-only' for='password'>Password</Label>
                <Input
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
                <Label class='sr-only' for='confirm_password'>Confirm Password</Label>
                <Input
                        id='confirm_password'
                        placeholder='Confirm Password'
                        type='password'
                        autocapitalize='none'
                        autocomplete='new-password'
                        autocorrect='off'
                        disabled={isLoading}
                        bind:value={confirm_password}
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
                <Checkbox id="terms" bind:checked={accepted_terms} class='border-input' />
                <Label for="terms" class="text-sm text-muted-foreground">I accept the <a href="#" class="text-primary">terms and conditions</a></Label>
            </div>

            <!-- Submit button -->
            <Button type='submit' disabled={isLoading} class='mt-2'>
                {#if isLoading}
                    <Icon icon='ph:spinner-bold' class='mr-2 h-4 w-4 animate-spin' />
                {/if}
                Submit
            </Button>

            <slot/>
        </div>
    </form>

</div>
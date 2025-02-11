import {user} from "$lib/stores";
import { goto } from '$app/navigation';
import { get } from 'svelte/store';
import {Sleep} from "$lib/utils";
import {throwToast} from "$lib/toasts";

const DEV_MODE = false;

type DelayedRedirect = {
    delay: number;
    title?: string;
    message?: string;
};

async function redirect(to: string, delayed?: DelayedRedirect) {
    if (delayed) {
        if (delayed.title) throwToast(delayed.title, delayed.message ?? '');
        await Sleep(delayed.delay);
    }

    console.log('[REDIRECT] Redirecting to', to);
    if (!DEV_MODE) goto(to);
}

async function redirectIfLoggedIn(to: string, delayed?: DelayedRedirect) {
    const store = get(user);
    if (!store.isLoggedIn) return;
    await redirect(to, delayed);
}

async function redirectIfLoggedOut(to: string, delayed?: DelayedRedirect) {
    const store = get(user);
    if (store.isLoggedIn) return;
    await redirect(to, delayed);
}

export { redirectIfLoggedIn, redirectIfLoggedOut, redirect };
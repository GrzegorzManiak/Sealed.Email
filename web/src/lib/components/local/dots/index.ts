import {writable} from "svelte/store";
import type {DotsSettings} from "$local/dots/index.d";

export { default as Dots } from './dots.svelte';

export const DefaultSettings: DotsSettings = {
    text: 'GRZEGORZ',
    text_height: 15,
    color: '#fff',
    dot_size: 2,
    dot_spacing: 1.5,
    force: 25,
    force_size: 50,
    max_dist: 150,
    data: {
        dots: [],
        rows: 0,
        cols: 0
    }
}

export const Small = {
    text_height: 15,
    dot_size: 2,
    dot_spacing: 1.5,
    width: 500,
};

export const Medium = {
    text_height: 16,
    dot_size: 5,
    dot_spacing: 4,
    width: 1000,
};

export const Large = {
    text_height: 18,
    dot_size: 7,
    dot_spacing: 6,
    width: 1500,
};

export const DotsWritable = writable<DotsSettings>(DefaultSettings);
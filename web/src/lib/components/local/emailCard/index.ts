import {DeviceType} from "$lib";

export { default as EmailCard } from "./emailCard.svelte";
export { default as ChainCard } from "./chainCard.svelte";
export * as Types from "./types";

const Sizes = {
    [DeviceType.Mobile]: {
        replyContactFontSize: '16px',
        contactFontSize: '18px',
        subjectFontSize: '16px',
        bodyFontSize: '14px',
        middleSectionWidth: '200px',
    },
    [DeviceType.Tablet]: {
        replyContactFontSize: '16px',
        contactFontSize: '18px',
        subjectFontSize: '16px',
        bodyFontSize: '14px',
        middleSectionWidth: '200px',
    },
    [DeviceType.Desktop]: {
        replyContactFontSize: '6px',
        contactFontSize: '18px',
        subjectFontSize: '16px',
        bodyFontSize: '14px',
        middleSectionWidth: '300px',
    },
};

export { Sizes };
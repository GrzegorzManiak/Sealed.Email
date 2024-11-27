import { writable } from "svelte/store";

enum DeviceType {
    Mobile = 'Mobile',
    Tablet = 'Tablet',
    Desktop = 'Desktop',
}

enum Sizes {
    Compact = 'Compact',
    Default = 'Default',
}

const ActiveDeviceType = writable<DeviceType>(DeviceType.Desktop);

export {
    DeviceType,
    ActiveDeviceType,
    Sizes
};
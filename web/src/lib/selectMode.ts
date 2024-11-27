import {writable, type Writable} from "svelte/store";

type ContextObject = {
    CurrentlySelected: Writable<string | null>;
    SelectMode: Writable<boolean>;
    SelectMap: Set<string>;
};

const selectMap: Map<string, ContextObject> = new Map();

function GetContextObject(group: string): ContextObject {
    let selectGroup = selectMap.get(group);
    if (!selectGroup) {
        selectGroup = {
            CurrentlySelected: writable<string | null>(null),
            SelectMode: writable<boolean>(false),
            SelectMap: new Set<string>()
        };
        selectMap.set(group, selectGroup);
    }
    return selectGroup;
}

function AddToSelectMap(id: string, group: string): void {
    const selectGroup = GetContextObject(group);
    selectGroup.SelectMap.add(id);
    selectGroup.SelectMode.set(selectGroup.SelectMap.size > 0);
    selectMap.set(group, selectGroup);
}

function RemoveFromSelectMap(id: string, group: string): void {
    const selectGroup = GetContextObject(group);
    selectGroup.SelectMap.delete(id);
    selectGroup.SelectMode.set(selectGroup.SelectMap.size > 0);
    selectMap.set(group, selectGroup);
}

function GetSelectMode(group: string): Writable<boolean> {
    const selectGroup = GetContextObject(group);
    return selectGroup.SelectMode;
}

function GetCurrentlySelected(group: string): Writable<string | null> {
    const selectGroup = GetContextObject(group);
    return selectGroup.CurrentlySelected;
}

function SetCurrentlySelected(group: string, id: string | null): void {
    const selectGroup = GetContextObject(group);
    selectGroup.CurrentlySelected.set(id);
    selectMap.set(group, selectGroup);
}

export {
    AddToSelectMap,
    RemoveFromSelectMap,
    SetCurrentlySelected,
    GetSelectMode,
    GetCurrentlySelected,
    selectMap
};
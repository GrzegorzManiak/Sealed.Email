import {writable} from "svelte/store";

type InboxContext = {
    selected: string | null;
    selectMode: boolean;
    chainIDS: Set<string>;
    grouped: Set<string>;
};

class InboxManager {

    private readonly id: string;
    private readonly name: string;

    private readonly grouped: Set<string> = new Set();
    private readonly currentlySelected: string | null = null;
    private readonly selectMode: boolean = false;
    private readonly chainIDS: Set<string> = new Set();

    private readonly inboxContext = writable<InboxContext>({
        selected: this.currentlySelected,
        selectMode: this.selectMode,
        chainIDS: this.chainIDS,
        grouped: this.grouped
    });

    public constructor(id: string, name: string) {
        this.id = id;
        this.name = name;
    }

    public chain(id: string): void {
        console.log(`Chain: ${id}`);
        this.inboxContext.update((context) => {
            if (context.chainIDS.has(id)) context.chainIDS.delete(id);
            else context.chainIDS.add(id);
            if (this.grouped.size > 0 && context.chainIDS.size === 0) context.selectMode = true;
            else context.selectMode = this.grouped.size > 1 || context.chainIDS.size > 1;
            return context;
        });
    }

    public select(id: string): void {
        console.log(`Select: ${id}`);
        this.inboxContext.update((context) => {
            if (context.selected === id) context.selected = null;
            else context.selected = id;
            return context;
        });
    }

    public check(id: string, value: boolean | null = null): void {
        console.log(`Check: ${id}`);
        this.inboxContext.update((context) => {
            if (value === null) {
                if (this.grouped.has(id)) this.grouped.delete(id);
                else this.grouped.add(id);
            }
            else {
                if (value) this.grouped.add(id);
                else this.grouped.delete(id);
            }
            if (this.grouped.size > 0 && context.chainIDS.size === 0) context.selectMode = true;
            else context.selectMode = this.grouped.size > 1 || context.chainIDS.size > 1;
            context.selected = null;
            return context;
        });
    }

    public resetCheck(): void {
        console.log("Reset Check");
        this.inboxContext.update((context) => {
            this.grouped.clear();
            context.selectMode = false;
            return context;
        });
    }

    public get rawContext() { return { selected: this.currentlySelected, selectMode: this.selectMode, chainIDS: this.chainIDS, grouped: this.grouped }; }
    public get context() { return this.inboxContext; }
}

export {
    InboxManager,
    type InboxContext
}
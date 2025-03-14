import Dexie from 'dexie';

abstract class StorageService {
	public abstract save(store: string, key: string, value: string): Promise<void>;
	public abstract get(store: string, key: string): Promise<string | null>;
	public abstract getDecryptionKey(): Promise<string | null>;
	public abstract setDecryptionKey(key: string): Promise<void>;
	public abstract reset(): Promise<void>;
}

class DummyStorageService extends StorageService {
	static store: Map<string, string> = new Map();

	public async save(store: string, key: string, value: string): Promise<void> {
		DummyStorageService.store.set(store + '-' + key, value);
	}

	public async get(store: string, key: string): Promise<string | null> {
		return DummyStorageService.store.get(store + '-' + key) ?? null;
	}

	public async getDecryptionKey(): Promise<string| null> {
		return this.get("key", "decryptionKey") ?? "";
	}

	public async setDecryptionKey(key: string): Promise<void> {
		this.save("key", "decryptionKey", key);
	}

	public async reset(): Promise<void> {
		DummyStorageService.store.clear();
	}
}

class IndexedDBStorageService extends StorageService {
	private db: Dexie;

	public constructor() {
		super();
		this.db = new Dexie('MyDatabase');
		this.db.version(1).stores({
			key: 'id'
		});
	}

	private async ensureStoreExists(store: string): Promise<void> {
		if (this.db.tables.some(table => table.name === store)) return;
		this.db.close();
		this.db.version(this.db.verno + 1).stores({[store]: 'id'});
		await this.db.open();
	}

	public async save(store: string, key: string, value: string): Promise<void> {
		await this.ensureStoreExists(store);
		await this.db.table(store).put({ id: key, value: value });
	}

	public async get(store: string, key: string): Promise<string | null> {
		await this.ensureStoreExists(store);
		const result = await this.db.table(store).get(key);
		return result ? result.value : null;
	}

	public async getDecryptionKey(): Promise<string | null> {
		return this.get("key", "decryptionKey");
	}

	public async setDecryptionKey(key: string): Promise<void> {
		return this.save("key", "decryptionKey", key);
	}

	public async reset(): Promise<void> {
		await this.ensureStoreExists("key");
		await this.db.table("key").clear();
	}
}


export {
	StorageService,
	DummyStorageService,
	IndexedDBStorageService
}
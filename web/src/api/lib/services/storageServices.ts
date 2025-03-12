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
	private dbName = "emailClientDB";
	private version = 1;

	private async getDB(storeName: string): Promise<IDBDatabase> {
		return new Promise((resolve, reject) => {
			const request = indexedDB.open(this.dbName, this.version);

			request.onupgradeneeded = (event) => {
				const db = (event.target as IDBOpenDBRequest).result;
				if (!db.objectStoreNames.contains(storeName)) {
					db.createObjectStore(storeName, { keyPath: "key" });
				}
			};

			request.onsuccess = () => resolve(request.result);
			request.onerror = () => reject(request.error);
		});
	}

	public async save(storeName: string, key: string, value: string): Promise<void> {
		const db = await this.getDB(storeName);
		return new Promise((resolve, reject) => {
			const tx = db.transaction(storeName, "readwrite");
			const store = tx.objectStore(storeName);
			const request = store.put({ key, value });

			request.onsuccess = () => resolve();
			request.onerror = () => reject(request.error);
		});
	}

	public async get(storeName: string, key: string): Promise<string | null> {
		const db = await this.getDB(storeName);
		return new Promise((resolve, reject) => {
			const tx = db.transaction(storeName, "readonly");
			const store = tx.objectStore(storeName);
			const request = store.get(key);

			request.onsuccess = () => resolve(request.result ? request.result.value : null);
			request.onerror = () => reject(request.error);
		});
	}

	public async getDecryptionKey(): Promise<string | null> {
		return this.get("info", "decryptionKey") ?? "";
	}

	public async setDecryptionKey(key: string): Promise<void> {
		this.save("info", "decryptionKey", key);
	}

	public async reset(): Promise<void> {
		const stores = ["info", "emailData", "emailMeta"];
		for (const store of stores) {
			const db = await this.getDB(store);
			const tx = db.transaction(store, "readwrite");
			tx.objectStore(store).clear();
		}
	}
}

export {
	StorageService,
	DummyStorageService
}
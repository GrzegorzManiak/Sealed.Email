import {StorageService} from "./storageServices";
import {Email} from "../api/email";
import Session from "../session/session";
import {UrlSafeBase64Decode, UrlSafeBase64Encode} from "../common";
import {Compress, Decompress, Decrypt, Encrypt, NewKey} from "../symetric";

type EmailMetadata = {
	emailID: string;
	receivedAt: number;
	bucketPath: string;
	read: boolean;
	spam: boolean;
	sent: boolean;
	folder: string;
	to: string;
	from: string;
	subject: string;
	storage: {
		encrypted: boolean;
		failed: boolean;
		local: boolean;
	}
};

class EmailStorage {
	private readonly storageService: StorageService;
	private readonly session: Session;

	private decryptionKey: Uint8Array | null = null;
	private ready: boolean = false;

	public constructor(storageService: StorageService, session: Session) {
		this.storageService = storageService;
		this.session = session;
	}

	public async init(): Promise<void> {
		if (this.ready) {
			console.log('EmailService already initialized');
			return;
		}

		const encodedEncryptionKey = await this.storageService.getDecryptionKey();
		if (!encodedEncryptionKey) await this.createEncryptionKey();
		else this.decryptionKey = UrlSafeBase64Decode(encodedEncryptionKey);
		await this.testKey();
		this.ready = true;

		console.log('EmailService initialized');
	}

	private async createEncryptionKey(): Promise<void> {
		this.storageService.reset();
		const key = NewKey();
		const encodedKey = UrlSafeBase64Encode(key);
		const encryptedKey = await this.session.EncryptKey(encodedKey);
		await this.storageService.setDecryptionKey(encryptedKey);
		this.decryptionKey = key;
		const sampleData = await Encrypt('sampleData', key);
		this.storageService.save('info', 'sampleData', UrlSafeBase64Encode(Compress(sampleData)));
	}

	private async testKey(): Promise<void> {
		if (!this.decryptionKey) throw new Error('No decryption key');
		const sampleData = await this.storageService.get('info', 'sampleData');

		if (!sampleData) {
			console.log('No sample data found, creating new key');
			return this.createEncryptionKey();
		}

		const decompressedData = Decompress(UrlSafeBase64Decode(sampleData));
		const decryptedData = await Decrypt(decompressedData, this.decryptionKey);

		if (decryptedData !== 'sampleData') {
			console.log('Decryption failed, creating new key');
			return this.createEncryptionKey();
		}
	}

	public async saveEmail(metaData: Email, content: string, subject: string): Promise<void> {
		if (!this.ready || !this.decryptionKey) throw new Error('EmailService not ready');
		const encryptedContent = await Encrypt(content, this.decryptionKey);
		const compressedContent = Compress(encryptedContent);

		const encryptedSubject = await Encrypt(subject, this.decryptionKey);
		const compressedSubject = Compress(encryptedSubject);

		await this.storageService.save('emailData', metaData.emailID, UrlSafeBase64Encode(compressedContent));
		await this.storageService.save('emailMeta', metaData.emailID, JSON.stringify({
			...metaData,
			subject: UrlSafeBase64Encode(compressedSubject),
			storage: { encrypted: true, failed: false, local: true }
		}));
	}

	public async updateMetaData(emailID: string, metaData: Email): Promise<void> {
		if (!this.ready || !this.decryptionKey) throw new Error('EmailService not ready');
		const existingMetaData = await this.getEmailMetadata(emailID);
		if (!existingMetaData) throw new Error('Email not found');
		await this.storageService.save('emailMeta', emailID, JSON.stringify({
			...existingMetaData,
			...metaData
		}));
	}

	public async getEmailMetadata(emailID: string): Promise<EmailMetadata | null> {
		if (!this.ready || !this.decryptionKey) throw new Error('EmailService not ready');
		const metaData = await this.storageService.get('emailMeta', emailID);
		if (!metaData) return null;
		return JSON.parse(metaData);
	}

	public async getEmailContent(emailID: string): Promise<string | null> {
		if (!this.ready || !this.decryptionKey) throw new Error('EmailService not ready');
		const compressedContent = await this.storageService.get('emailData', emailID);
		if (!compressedContent) return null;
		const decompressedContent = Decompress(UrlSafeBase64Decode(compressedContent));
		return Decrypt(decompressedContent, this.decryptionKey);
	}

	public get isReady(): boolean {
		return this.ready;
	}
}

export default EmailStorage;

export {
	type EmailMetadata
}
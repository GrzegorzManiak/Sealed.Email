import {StorageService} from "./storageServices";
import Session from "../session/session";
import {UrlSafeBase64Decode, UrlSafeBase64Encode} from "../common";
import {Compress, Decompress, Decrypt, Encrypt, NewKey} from "../symetric";

type Address = {
	address?: string;
	name: string;
}

type EmailMetadata = {
	emailID: string;
	receivedAt: number;
	bucketPath: string;
	read: boolean;
	spam: boolean;
	sent: boolean;
	folder: string;
	to: Array<Address>;
	cc: Array<Address>;
	from: Address;
	replyTo: Array<Address>;
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
		if (this.ready) return

		const encodedEncryptionKey = await this.storageService.getDecryptionKey();
		if (!encodedEncryptionKey) await this.createEncryptionKey();
		else try {
			const decryptedKey = await this.session.DecryptKey(encodedEncryptionKey);
			this.decryptionKey = UrlSafeBase64Decode(decryptedKey);
		} catch (error) { this.decryptionKey = await this.createEncryptionKey(); }

		await this.testKey();
		this.ready = true;
	}

	private async createEncryptionKey(): Promise<Uint8Array> {
		await this.storageService.reset();
		const key = NewKey(32);
		const encryptedKey = await this.session.EncryptKey(key);
		await this.storageService.setDecryptionKey(encryptedKey);
		this.decryptionKey = key;
		const sampleData = await Encrypt('sampleData', key);
		await this.storageService.save('info', 'sampleData', UrlSafeBase64Encode(Compress(sampleData)));
		return key;
	}

	private async testKey(): Promise<Uint8Array> {
		if (!this.decryptionKey) throw new Error('No decryption key');
		const sampleData = await this.storageService.get('info', 'sampleData');
		if (!sampleData) return await this.createEncryptionKey();
		const decompressedData = Decompress(UrlSafeBase64Decode(sampleData));
		try {
			const decryptedData = await Decrypt(decompressedData, this.decryptionKey);
			if (decryptedData !== 'sampleData') return await this.createEncryptionKey();
		} catch (error) { return await this.createEncryptionKey(); }
		return this.decryptionKey;
	}

	public async saveEmail(metaData: EmailMetadata, content: string): Promise<void> {
		if (!this.ready || !this.decryptionKey) throw new Error('EmailService not ready');

		const encryptedContent = await Encrypt(content, this.decryptionKey);
		const compressedContent = Compress(encryptedContent);

		const encryptedSubject = await Encrypt(metaData.subject, this.decryptionKey);
		const compressedSubject = Compress(encryptedSubject);

		metaData.subject = UrlSafeBase64Encode(compressedSubject);

		await this.storageService.save('emailData', metaData.emailID, UrlSafeBase64Encode(compressedContent));
		await this.storageService.save('emailMeta', metaData.emailID, JSON.stringify({
			...metaData,
			storage: { encrypted: true, failed: false, local: true }
		}));
	}

	public async updateMetaData(emailID: string, metaData: Partial<EmailMetadata>): Promise<void> {
		if (!this.ready || !this.decryptionKey) throw new Error('EmailService not ready');
		const existingMetaData = await this.getEmailMetadata(emailID, false);
		if (!existingMetaData) throw new Error('Email not found')
		const updatedMetaData = {
			...existingMetaData,
			spam: metaData.spam ?? existingMetaData.spam,
			read: metaData.read ?? existingMetaData.read,
			folder: metaData.folder ?? existingMetaData.folder,
		}

		await this.storageService.save('emailMeta', emailID, JSON.stringify(updatedMetaData));
	}

	public async getEmailMetadata(emailID: string, decryptSubject: boolean = true): Promise<EmailMetadata | null> {
		if (!this.ready || !this.decryptionKey) throw new Error('EmailService not ready');
		const metaData = await this.storageService.get('emailMeta', emailID);
		if (!metaData) return null;

		const parsedMetaData = JSON.parse(metaData) as EmailMetadata;
		if (decryptSubject) {
			const compressedSubject = UrlSafeBase64Decode(parsedMetaData.subject);
			parsedMetaData.subject = await Decrypt(Decompress(compressedSubject), this.decryptionKey);
		}
		return parsedMetaData;
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
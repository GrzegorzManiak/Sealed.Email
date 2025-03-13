import {EmailStorage, Session} from "../index";
import Domain from "./domain";
import {Email, EmailListFilters, GetEmailData, GetEmailList} from "../api/email";
import {EmailMetadata} from "./emailStorage";
import PostalMime, {type Email as ParsedEmail, Header} from 'postal-mime';
import { extract } from 'letterparser';
import {UrlSafeBase64Decode, UrlSafeBase64Encode} from "../common";
import { Decrypt as AsymDecrypt, ExtractEphemeralPubKey } from "../asymmetric";
import { Decompress, Decrypt } from "../symetric";

const PostEncryptionHeader = 'x-noise-post-encryption-keys';
const InboxKeysHeader = 'x-noise-inbox-keys';

// TODO: Private key version

type PostEncryptionRecipient = { publicKey: string, emailKey: string }


function FormatPostEncryptionHeader(rawRecord: string): PostEncryptionRecipient | null {
	rawRecord = rawRecord.replace(/\s/g, '');
	const parts = rawRecord.split('<', 2);
	if (parts.length !== 2) return null;
	const publicKey = parts[0];
	const emailKey = parts[1].slice(0, -1);
	return { publicKey, emailKey };
}

function DecodePostEncryptionHeader(header: string, domain: Domain): PostEncryptionRecipient | null {
	const records = header.split(',');
	for (const record of records) {
		const recipient = FormatPostEncryptionHeader(record);
		if (!recipient) continue;
		if (domain.HasPublicKey(recipient.publicKey)) return recipient;
	}
	return null;
}

class EmailProvider {
	private readonly _emailStorage: EmailStorage;
	private readonly _session: Session;

	public constructor(
		emailStorage: EmailStorage,
		session: Session
	) {
		if (!emailStorage.isReady) throw new Error('EmailStorage not ready');
		this._emailStorage = emailStorage;
		this._session = session;
	}

	private async decryptPostEncryptionEmail(postEncryptionKeys: Header, domain: Domain, parsedEmail: ParsedEmail, emailData: string): Promise<string> {
		const keys = DecodePostEncryptionHeader(postEncryptionKeys.value, domain);
		if (!keys) throw new Error('No valid post encryption keys found');

		const privateKey = domain.GetPrivateKey(keys.publicKey);
		if (!privateKey) throw new Error('No private key found');

		const decodedEmailKey = UrlSafeBase64Decode(keys.emailKey);
		const decodedData = ExtractEphemeralPubKey(decodedEmailKey);
		const decryptedEmailKey = await AsymDecrypt(decodedData.ephemeralPub, privateKey, decodedData.ciphertext);

		const body = parsedEmail.attachments[0];
		if (body.mimeType !== 'application/json' || !body) throw new Error('Invalid body type');
		const rawContent = emailData.split('\r\n\r\n', 2)[1].split('\r\n').join('');
		if (!rawContent) throw new Error('Invalid email data');

		const decompressedContent = Decompress(UrlSafeBase64Decode(rawContent));
		return await Decrypt(decompressedContent, UrlSafeBase64Decode(decryptedEmailKey));
	}

	private async decryptInboxKeys(header: Header, domain: Domain): Promise<string | null> {
		let value = header.value;
		value = value.replace(/\s/g, '');
		const keys = value.split(',');

		for (const key of keys) {
			const parts = key.split(':', 2);
			parts[0] = parts[0].replace('<', '');
			parts[1] = parts[1].replace('>', '');

			const privateKey = domain.GetPrivateKey(parts[0]);
			if (!privateKey) continue;

			const decodedEmailKey = UrlSafeBase64Decode(parts[1]);
			const decodedData = ExtractEphemeralPubKey(decodedEmailKey);
			return AsymDecrypt(decodedData.ephemeralPub, privateKey, decodedData.ciphertext);
		}

		return null;
	}

	private async decryptCc(email: ParsedEmail, emailKey: Uint8Array) {
		if (!email.cc || email.cc.length === 0) return;
		for (let i = 0; i < email.cc.length; i++) {
			const cc = email.cc![i];
			email.cc![i].name = await Decrypt(Decompress(UrlSafeBase64Decode(cc.name)), emailKey);
		}
	}

	private async decryptTo(email: ParsedEmail, emailKey: Uint8Array) {
		if (!email.to || email.to.length === 0) return;
		for (let i = 0; i < email.to.length; i++) {
			const to = email.to![i];
			email.to![i].name = await Decrypt(Decompress(UrlSafeBase64Decode(to.name)), emailKey);
		}
	}

	private async decryptFrom(email: ParsedEmail, emailKey: Uint8Array) {
		email.from.name = await Decrypt(Decompress(UrlSafeBase64Decode(email.from.name)), emailKey);
	}

	private async decryptReplyTo(email: ParsedEmail, emailKey: Uint8Array) {
		if (!email.replyTo || email.replyTo.length === 0) return;
		for (let i = 0; i < email.replyTo.length; i++) {
			const replyTo = email.replyTo![i];
			email.replyTo![i].name = await Decrypt(Decompress(UrlSafeBase64Decode(replyTo.name)), emailKey);
		}
	}

	private async decryptSubject(email: ParsedEmail, emailKey: Uint8Array) {
		if (!email.subject) return;
		email.subject = atob(await Decrypt(Decompress(UrlSafeBase64Decode(email.subject)), emailKey));
	}

	private async decryptText(email: string, emailKey: Uint8Array): Promise<string> {
		const rawContent = email.split('\r\n\r\n', 2)[1].split('\r\n').join('');
		if (!rawContent) throw new Error('Invalid email data');
		return atob(await Decrypt(Decompress(UrlSafeBase64Decode(rawContent)), emailKey));
	}

	private async pullEmailData(emailID: string, email: Email, domain: Domain): Promise<EmailMetadata> {
		const emailData = await GetEmailData(this._session, emailID, email);
		if (!emailData) throw new Error('Email not found');

		const parsedEmail = await PostalMime.parse(emailData);
		let decryptedContent: ParsedEmail = parsedEmail;
		let metadata: EmailMetadata;
		let emailContent: string;

		const postEncryptionKeys = parsedEmail.headers.find(header => header.key === PostEncryptionHeader);
		if (postEncryptionKeys !== undefined) {
			const decrypted = await this.decryptPostEncryptionEmail(postEncryptionKeys, domain, parsedEmail, emailData);
			decryptedContent = await PostalMime.parse(decrypted);
			emailContent = decrypted.split('\r\n\r\n', 2)[1].split('\r\n').join('');
			metadata = {
				emailID,
				from: decryptedContent.from,
				cc: decryptedContent.cc ?? [],
				to: decryptedContent.to ?? [],
				replyTo: decryptedContent.replyTo ?? [decryptedContent.from],

				subject: decryptedContent.subject ?? '',
				bucketPath: email.bucketPath,
				folder: email.folder,
				read: email.read,
				sent: email.sent,
				spam: email.spam,
				receivedAt: email.receivedAt,

				storage: {
					encrypted: false,
					local: true,
					failed: false,
				}
			}
		}

		else {
			const sharedKeys = parsedEmail.headers.find(header => header.key === InboxKeysHeader);
			if (sharedKeys === undefined) throw new Error('No shared keys found');
			const decryptedKeys = await this.decryptInboxKeys(sharedKeys, domain);
			if (!decryptedKeys) throw new Error('Failed to decrypt shared keys');

			const emailKey = UrlSafeBase64Decode(decryptedKeys);
			await this.decryptCc(parsedEmail, emailKey);
			await this.decryptTo(parsedEmail, emailKey);
			await this.decryptFrom(parsedEmail, emailKey);
			await this.decryptReplyTo(parsedEmail, emailKey);
			await this.decryptSubject(parsedEmail, emailKey);
			emailContent = await this.decryptText(emailData, emailKey);

			metadata = {
				emailID,
				from: parsedEmail.from,
				cc: parsedEmail.cc ?? [],
				to: parsedEmail.to ?? [],
				replyTo: parsedEmail.replyTo ?? [parsedEmail.from],

				subject: parsedEmail.subject ?? '',
				bucketPath: email.bucketPath,
				folder: email.folder,
				read: email.read,
				sent: email.sent,
				spam: email.spam,
				receivedAt: email.receivedAt,

				storage: {
					encrypted: true,
					local: true,
					failed: false,
				}
			}
		}

		await this._emailStorage.saveEmail(
			metadata,
			emailContent
		);

		return metadata;
	}

	public async getEmails(domain: Domain, filter: EmailListFilters): Promise<Array<EmailMetadata>> {
		const emails = await GetEmailList(this._session, filter);
		const emailList: Array<Promise<EmailMetadata> | EmailMetadata> = [];
		console.log(emails);
		for (const email of emails.emails) {

			try {
				await this._emailStorage.updateMetaData(email.emailID, {
					spam: email.spam,
					read: email.read,
					folder: email.folder,
				});
			}

			catch (error) {
				console.warn('Failed to update email metadata:', error);
			}

			const metadata = await this._emailStorage.getEmailMetadata(email.emailID);
			if (metadata) {
				emailList.push(metadata);
				continue;
			}

			emailList.push(this.pullEmailData(email.emailID, email, domain));
		}

		return Promise.all(emailList);
	}
}

export default EmailProvider;
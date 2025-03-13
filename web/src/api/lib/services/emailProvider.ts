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

	private async pullEmailData(emailID: string, email: Email, domain: Domain): Promise<EmailMetadata> {
		const emailData = await GetEmailData(this._session, emailID, email);
		if (!emailData) throw new Error('Email not found');

		const parsedEmail = await PostalMime.parse(emailData);
		let decryptedContent: ParsedEmail = parsedEmail;

		const postEncryptionKeys = parsedEmail.headers.find(header => header.key === PostEncryptionHeader);
		if (postEncryptionKeys !== undefined) {
			const decrypted = await this.decryptPostEncryptionEmail(postEncryptionKeys, domain, parsedEmail, emailData);
			decryptedContent = await PostalMime.parse(decrypted);
		}

		console.log(decryptedContent);
	}

	public async getEmails(domain: Domain, filter: EmailListFilters): Promise<Array<EmailMetadata>> {
		const emails = await GetEmailList(this._session, filter);
		const emailList: Array<Promise<EmailMetadata> | EmailMetadata> = [];

		for (const email of emails.emails) {
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
import Domain from "./domain";
import {EncryptedInbox, Session} from "../index";
import * as Sym from "../symetric";
import {Hash} from "gowl-client-lib";
import {PostEncryptedEmail, SendEncryptedEmail} from "../api/email";
import {UrlSafeBase64Encode} from "../common";

type EncryptedEmailContent = {
	domain: Domain;
	key: Uint8Array;
	from: EncryptedInbox;
	to: EncryptedInbox;
	subject: string;
	body: string;
	cc?: EncryptedInbox[];
	bcc?: EncryptedInbox[];
	inReplyTo?: string;
	references?: string[];
};

class EncryptedEmail {
	private readonly _domain: Domain;
	private readonly _key: Uint8Array;

	private readonly _from: EncryptedInbox;
	private readonly _to: EncryptedInbox;
	private readonly _cc: EncryptedInbox[];
	private readonly _bcc: EncryptedInbox[];

	private readonly _subject: string;
	private readonly _body: string;

	private readonly _inReplyTo: string;
	private readonly _references: string[];
	private _signature: string | null = null;

	public constructor(content: EncryptedEmailContent) {
		this._domain = content.domain;
		this._from = content.from;
		this._to = content.to;
		this._cc = content.cc || [];
		this._bcc = content.bcc || [];
		this._subject = content.subject;
		this._body = content.body;
		this._inReplyTo = content.inReplyTo || "";
		this._references = content.references || [];
		this._key = content.key;
	}

	public async EncryptBody(): Promise<string> {
		const b64Body = Buffer.from(this._body).toString('base64');
		const encryptedBody = await Sym.Encrypt(b64Body, this._key);
		const compressedBody = Sym.Compress(encryptedBody);
		return UrlSafeBase64Encode(compressedBody);
	}

	public async EncryptSubject(): Promise<string> {
		const b64Subject = Buffer.from(this._subject).toString('base64');
		const encryptedSubject = await Sym.Encrypt(b64Subject, this._key);
		const compressedSubject = Sym.Compress(encryptedSubject);
		return UrlSafeBase64Encode(compressedSubject);
	}

	public async Encrypt(): Promise<PostEncryptedEmail> {
		const body = await this.EncryptBody();
		const subject = await this.EncryptSubject();

		return {
			domainID: this._domain.DomainID,
			subject,
			body,
			from: this._from.ComputedEncryptedInbox,
			to: this._to.ComputedEncryptedInbox,
			cc: this._cc.map((inbox) => inbox.ComputedEncryptedInbox),
			bcc: this._bcc.map((inbox) => inbox.ComputedEncryptedInbox),
			inReplyTo: this._inReplyTo,
			references: this._references,
		};
	}

	public async Sign(): Promise<string> {
		let data = this._body;
		data += '\n' + this._subject;
		data += '\n' + this._from.ComputedStringifiedInbox;
		data += '\n' + this._to.ComputedStringifiedInbox;

		const cc = this._cc.map((inbox) => inbox.ComputedStringifiedInbox);
		cc.sort((a, b) => a.localeCompare(b));

		data += '\n' + cc.join('\n');
		data += '\n' + this._inReplyTo;
		data += '\n' + this._references.join('\n');

		const dataHash = await Hash(data);
		this._signature = await this._domain.SignData(UrlSafeBase64Encode(dataHash));
		return this._signature;
	}

	public async Send(session: Session): Promise<void> {
		if (!this._signature) await this.Sign();
		const email = await this.Encrypt();
		if (!this._signature) throw new Error("Failed to sign email");
		await SendEncryptedEmail(session, email, this._signature);
	}
}

export default EncryptedEmail;

export {
	EncryptedEmailContent,
}
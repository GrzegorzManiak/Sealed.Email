import Domain from "./domain";
import {EncryptedInbox} from "../index";

class EncryptedEmail {
	private readonly _domain: Domain;

	private readonly _from: EncryptedInbox;
	private readonly _to: EncryptedInbox;
	private readonly _cc: EncryptedInbox[];
	private readonly _bcc: EncryptedInbox[];

	private readonly _subject: string;
	private readonly _body: string;

	private readonly _inReplyTo: string;
	private readonly _references: string[];
	private readonly _signature: string;

	public constructor(
		domain: Domain,
		from: EncryptedInbox,
		to: EncryptedInbox,
		subject: string,
		body: string,
		cc: EncryptedInbox[] = [],
		bcc: EncryptedInbox[] = [],
		inReplyTo: string = "",
		references: string[] = []
	) {
		this._domain = domain;
		this._from = from;
		this._to = to;
		this._cc = cc;
		this._bcc = bcc;
		this._subject = subject;
		this._body = body;
		this._inReplyTo = inReplyTo;
		this._references = references;
		this._signature = "";
	}
}

export default EncryptedEmail;
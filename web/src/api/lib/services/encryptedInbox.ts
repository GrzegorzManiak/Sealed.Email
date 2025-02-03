import Domain from "./domain";
import * as Sym from "../symetric";
import * as Asym from "../asymmetric";
import {DecodeFromBase64} from "../common";
import {EncodeToBase64} from "gowl-client-lib";

class EncryptedInbox {
	private readonly _domain: Domain;
	private readonly _email: string;
	private readonly _displayName: string;
	private readonly _publicKey: string;
	private readonly _emailKey: string;
	private readonly _encryptedEmailKey: string;

	public constructor(
		domain: Domain,
		email: string,
		displayName: string,
		publicKey: string,
		emailKey: string,
		encryptedEmailKey: string
	) {
		this._domain = domain;
		this._email = email;
		this._displayName = displayName;
		this._publicKey = publicKey;
		this._emailKey = emailKey;
		this._encryptedEmailKey = encryptedEmailKey;
	}

	public static async Create(
		domain: Domain,
		email: string,
		displayName: string,
		publicKey: string,
		emailKey: string
	): Promise<EncryptedInbox> {
		const encryptedName = Sym.Compress(await Sym.Encrypt(displayName, DecodeFromBase64(emailKey)));
		const encryptedEmailKey = await Asym.Encrypt(emailKey, DecodeFromBase64(publicKey));
		return new EncryptedInbox(domain, email, EncodeToBase64(encryptedName), publicKey, emailKey, encryptedEmailKey);
	}

	public get domain(): Domain {
		return this._domain;
	}

	public get email(): string {
		return this._email;
	}

	public get displayName(): string {
		return this._displayName;
	}
}

export default EncryptedInbox;
import Domain from "./domain";
import * as Sym from "../symetric";
import * as Asym from "../asymmetric";
import {DecodeFromBase64, EnsureFqdn, SplitEmail} from "../common";
import {EncodeToBase64, Hash} from "gowl-client-lib";
import {ComputedEncryptedInbox} from "../api/email";

class EncryptedInbox {
	private readonly _email: string;
	private readonly _displayName: string;
	private readonly _encryptedDisplayName: string;
	private readonly _publicKey: Uint8Array;
	private readonly _emailKey: Uint8Array;
	private readonly _encryptedEmailKey: string;
	private readonly _userHash: string;

	public constructor(
		email: string,
		displayName: string,
		publicKey: Uint8Array,
		emailKey: Uint8Array,
		encryptedEmailKey: string,
		userHash: string,
		encryptedDisplayName: string
	) {
		this._email = email;
		this._displayName = displayName;
		this._publicKey = publicKey;
		this._emailKey = emailKey;
		this._encryptedEmailKey = encryptedEmailKey;
		this._userHash = userHash;
		this._encryptedDisplayName = encryptedDisplayName;
	}

	public static async Create(
		email: string,
		displayName: string,
		publicKey: Uint8Array,
		emailKey: Uint8Array
	): Promise<EncryptedInbox> {

		const encryptedDisplayName = Sym.Compress(await Sym.Encrypt(displayName, emailKey));
		const encodedDisplayName = EncodeToBase64(encryptedDisplayName);

		const { domain, username } = SplitEmail(email);
		if (!domain || !username) throw new Error("Invalid email");
		const userHash = await Hash(`${username}@${EnsureFqdn(domain)}`);
		const encodedUserHash = `${EncodeToBase64(userHash)}@${domain}`;

		const encryptedEmailKey = await Asym.Encrypt(emailKey, publicKey);

		return new EncryptedInbox(
			email,
			displayName,
			publicKey,
			emailKey,
			encryptedEmailKey,
			encodedUserHash,
			encodedDisplayName
		);
	}

	public get ComputedStringifiedInbox(): string {
		const computedInbox = this.ComputedEncryptedInbox;
		return `${computedInbox.displayName}.${computedInbox.emailHash}.${computedInbox.publicKey}.${computedInbox.encryptedEmailKey}`;
	}

	public get ComputedEncryptedInbox(): ComputedEncryptedInbox {
		return {
			displayName: this._encryptedDisplayName,
			emailHash: this._userHash,
			publicKey: EncodeToBase64(this._publicKey),
			encryptedEmailKey: this._encryptedEmailKey
		}
	}
}

export default EncryptedInbox;
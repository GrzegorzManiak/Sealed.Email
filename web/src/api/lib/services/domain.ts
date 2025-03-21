
import Session from "../session/session";
import {DomainDnsData, DomainFull, DomainRefID} from "../api/domain";
import {Compress, Decompress, Decrypt, Encrypt} from "../symetric";
import { SignData} from "../asymmetric";
import {UrlSafeBase64Decode, UrlSafeBase64Encode} from "../common";
import EncryptedInbox from "./encryptedInbox";

class Domain {
    private readonly _domainRaw: DomainFull;
    private readonly _domainID: DomainRefID;
    private readonly _domain: string;
    private readonly _dateAdded: Date;

    private _symmetricRootKey: string;
    private _decryptedRootKey: Uint8Array;
    private _publicKey: string;
    private _privateKey: Uint8Array;

    private _dns: DomainDnsData;
    private _version: number;
    private _verified: boolean;
    private _catchAll: boolean;

    public constructor(
        domain: DomainFull,
        decryptedRootKey: Uint8Array,
        privateKey: Uint8Array,
    ) {
        this._domainRaw = domain;
        this._domainID = domain.domainID;
        this._domain = domain.domain;
        this._dateAdded = new Date(domain.dateAdded);

        this._symmetricRootKey = domain.symmetricRootKey;
        this._decryptedRootKey = decryptedRootKey;
        this._publicKey = domain.publicKey;
        this._privateKey = privateKey;

        this._dns = domain.dns;
        this._version = domain.version;
        this._verified = domain.verified;
        this._catchAll = domain.catchAll;

    }

    public static async Decrypt(
        session: Session,
        domain: DomainFull
    ): Promise<Domain> {
        const encodedRootKey = await session.DecryptKey(domain.symmetricRootKey);
        const rootKey = UrlSafeBase64Decode(encodedRootKey);

        const decodedPrivateKey = UrlSafeBase64Decode(domain.encryptedPrivateKey);
        const decompressedPrivateKey = Decompress(decodedPrivateKey);
        const encodedPrivateKey = await Decrypt(decompressedPrivateKey, rootKey);
        const privateKey = UrlSafeBase64Decode(encodedPrivateKey);

        return new Domain(domain, rootKey, privateKey);
    }

    public async EncryptData(key: Uint8Array | string, privateKey: Uint8Array = this._decryptedRootKey): Promise<string> {
        if (typeof key !== 'string') key = UrlSafeBase64Encode(key);
        const encryptedKey = Compress(await Encrypt(key, privateKey));
        return UrlSafeBase64Encode(encryptedKey);
    }

    public async DecryptData(key: Uint8Array | string, privateKey: Uint8Array = this._decryptedRootKey): Promise<string> {
        if (typeof key !== 'string') key = UrlSafeBase64Encode(key);
        const decompressedKey = Decompress(UrlSafeBase64Decode(key));
        return await Decrypt(decompressedKey, privateKey);
    }

    public async SignData(data: string, privateKey: Uint8Array = this._decryptedRootKey): Promise<string> {
        return await SignData(data, privateKey);
    }

    public FormatEmail(user: string): string {
        user = user.toLowerCase();
        user.replace(/[^a-z0-9]/g, '');
        return `${user}@${this._domain}`;
    }
    
    public async GetSender(emailKey: Uint8Array, user: string, displayName: string = ''): Promise<EncryptedInbox> {
        return EncryptedInbox.Create(
            this.FormatEmail(user),
            displayName,
            UrlSafeBase64Decode(this._publicKey),
            emailKey
        );
    }

    public HasPublicKey(publicKey: string): boolean {
        // TODO: This will have a versioning support later
        return this._publicKey === publicKey;
    }

    public GetPrivateKey(publicKey: string): Uint8Array | null {
        // TODO: This will have a versioning support later
        if (this._publicKey !== publicKey) return null;
        return this._privateKey;
    }

    public get DomainID(): DomainRefID { return this._domainID; }
    public get Domain(): string { return this._domain; }
    public get PublicKey(): string { return this._publicKey; }
}

export default Domain;
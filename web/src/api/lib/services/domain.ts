import Session from "../session/session";
import {DomainDnsData, DomainFull, DomainRefID} from "../api/domain";
import {Compress, Decompress, Decrypt, Encrypt, NewKey} from "../symetric";
import {BigIntToByteArray, EncodeToBase64, GetCurve, Hash} from "gowl-client-lib";
import {CurrentCurve} from "../constants";
import {DecodeFromBase64} from "../common";
import {InboxKeys} from "../api/inbox";

class Domain {
    private readonly _session: Session;
    private readonly _domainRaw: DomainFull;
    private readonly _domainID: DomainRefID;
    private readonly _domain: string;
    private readonly _dateAdded: Date;

    private _symmetricRootKey: string;
    private _decryptedRootKey: Uint8Array;

    private _dns: DomainDnsData;
    private _version: number;
    private _verified: boolean;
    private _catchAll: boolean;

    public constructor(
        session: Session,
        domain: DomainFull,
        decryptedRootKey: Uint8Array
    ) {
        this._session = session;
        this._domainRaw = domain;
        this._domainID = domain.domainID;
        this._domain = domain.domain;
        this._dateAdded = new Date(domain.dateAdded);

        this._symmetricRootKey = domain.symmetricRootKey;
        this._decryptedRootKey = decryptedRootKey;

        this._dns = domain.dns;
        this._version = domain.version;
        this._verified = domain.verified;
        this._catchAll = domain.catchAll;
    }

    public static async Decrypt(
        session: Session,
        domain: DomainFull
    ): Promise<Domain> {
        const rootKey = await session.DecryptKey(domain.symmetricRootKey);
        return new Domain(session, domain, DecodeFromBase64(rootKey));
    }

    public async CreateEmailHash(email: string): Promise<string> {
        const hash = BigIntToByteArray(await Hash(`${email}@${this._domain}`));
        return EncodeToBase64(hash);
    }

    public async CreateInboxKeys(inboxName: string): Promise<InboxKeys> {
        const symmetricRootKey = NewKey();
        const curve = GetCurve(CurrentCurve);
        const asymmetricPrivateKey = curve.utils.randomPrivateKey();
        const asymmetricPublicKey = curve.getPublicKey(asymmetricPrivateKey);
        const encryptedSymmetricRootKey = await this.EncryptKey(symmetricRootKey);
        const encryptedAsymmetricPrivateKey = await Encrypt(EncodeToBase64(asymmetricPrivateKey), symmetricRootKey);
        const encryptedInboxName = await this.EncryptKey(inboxName);

        return {
            emailHash: await this.CreateEmailHash(inboxName),
            symmetricRootKey: encryptedSymmetricRootKey,
            asymmetricPrivateKey: EncodeToBase64(Compress(encryptedAsymmetricPrivateKey)),
            asymmetricPublicKey: EncodeToBase64(asymmetricPublicKey),
            encryptedInboxName: encryptedInboxName
        }
    }

    public async EncryptKey(key: Uint8Array | string): Promise<string> {
        if (typeof key !== 'string') key = EncodeToBase64(key);
        const encryptedKey = Compress(await Encrypt(key, this._decryptedRootKey));
        return EncodeToBase64(encryptedKey);
    }

    public async DecryptKey(key: Uint8Array | string): Promise<string> {
        if (typeof key !== 'string') key = EncodeToBase64(key);
        const decompressedKey = Decompress(DecodeFromBase64(key));
        return await Decrypt(decompressedKey, this._decryptedRootKey);
    }

    public get DomainID(): DomainRefID { return this._domainID; }
    public get Domain(): string { return this._domain; }
}

export default Domain;
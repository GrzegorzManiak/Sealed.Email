import {Statistics} from "./types";
import {ReturnedVerifyData} from "../auth/types";
import { Login } from "../auth/login";
import {Compress, Decompress, Decrypt, Encrypt} from "../symetric";
import {UrlSafeBase64Decode, UrlSafeBase64Encode} from "../common";
import {CryptoGenericError} from "../errors";
import {StandardIntegrityHash} from "../auth/register";
import {BigIntToByteArray} from "gowl-client-lib";
import {COOKIE_NAME} from "../constants";

interface SerializedSession {
    encryptedSymetricRootKey: string;
    encryptedAsymetricPrivateKey: string;
    encryptedSymetricContactsKey: string;
    integrityHash: string;
    sessionEstablished: boolean;
    sessionToken: string;
    statistics: Statistics;

    passwordHash: string;
    rootKey: string;
    privateKey: string;
    contactsKey: string;
    sessionKey: string;
}

class Session {
    public static Deserialize(data: string): Session | Error {
        try {
            const parsedData = JSON.parse(data) as SerializedSession;
            const session = new Session(null as any);
            session._encryptedSymetricRootKey = parsedData.encryptedSymetricRootKey;
            session._encryptedAsymetricPrivateKey = parsedData.encryptedAsymetricPrivateKey;
            session._encryptedSymetricContactsKey = parsedData.encryptedSymetricContactsKey;
            session._integrityHash = parsedData.integrityHash;
            session._sessionEstablished = parsedData.sessionEstablished;
            session._sessionToken = parsedData.sessionToken;
            session._statistics = parsedData.statistics;

            session._passwordHash = UrlSafeBase64Decode(parsedData.passwordHash);
            session._rootKey = UrlSafeBase64Decode(parsedData.rootKey);
            session._privateKey = UrlSafeBase64Decode(parsedData.privateKey);
            session._contactsKey = UrlSafeBase64Decode(parsedData.contactsKey);
            session._sessionKey = UrlSafeBase64Decode(parsedData.sessionKey);

            return session;
        } catch (error) {
            return new Error(`Failed to deserialize session: ${error}`);
        }
    }

    private _encryptedSymetricRootKey: string;
    private _encryptedAsymetricPrivateKey: string;
    private _encryptedSymetricContactsKey: string;
    private _integrityHash: string;

    private _sessionEstablished: boolean = false;
    private _sessionToken: string = '';

    private _passwordHash: Uint8Array;

    private _rootKey: Uint8Array;
    private _privateKey: Uint8Array;
    private _contactsKey: Uint8Array;
    private _sessionKey: Uint8Array;

    private _statistics: Statistics;

    public constructor(awaitedLogin: Awaited<ReturnType<typeof Login>>, captureCookie: boolean = false) {
        if (awaitedLogin === null) {
            this._encryptedSymetricRootKey = '';
            this._encryptedAsymetricPrivateKey = '';
            this._encryptedSymetricContactsKey = '';
            this._integrityHash = '';
            this._passwordHash = new Uint8Array();
            this._rootKey = new Uint8Array();
            this._privateKey = new Uint8Array();
            this._contactsKey = new Uint8Array();
            this._sessionKey = new Uint8Array();
            this._statistics = {
                totalInboundEmails: 0,
                totalInboundBytes: 0,
                totalOutboundEmails: 0,
                totalOutboundBytes: 0
            };
            return;
        }

        const { verify, passwordHash, client } = awaitedLogin;
        this._passwordHash = passwordHash;

        const sessionKey = client.GetSessionKey();
        if (sessionKey instanceof Error) throw sessionKey;
        this._sessionKey = BigIntToByteArray(sessionKey);
        if (captureCookie) this._sessionToken = this._FindSessionToken(verify._headers);

        this._encryptedSymetricRootKey = verify.encryptedSymmetricRootKey;
        this._encryptedAsymetricPrivateKey = verify.encryptedAsymmetricPrivateKey;
        this._encryptedSymetricContactsKey = verify.encryptedSymmetricContactsKey;
        this._integrityHash = verify.integrityHash;

        this._rootKey = new Uint8Array();
        this._privateKey = new Uint8Array();
        this._contactsKey = new Uint8Array();

        this._statistics = {
            totalInboundEmails: verify.totalInboundEmails,
            totalInboundBytes: verify.totalInboundBytes,
            totalOutboundEmails: verify.totalOutboundEmails,
            totalOutboundBytes: verify.totalOutboundBytes
        };
    }

    public Serialize(): string {
        const serializedData: SerializedSession = {
            encryptedSymetricRootKey: this._encryptedSymetricRootKey,
            encryptedAsymetricPrivateKey: this._encryptedAsymetricPrivateKey,
            encryptedSymetricContactsKey: this._encryptedSymetricContactsKey,
            integrityHash: this._integrityHash,
            sessionEstablished: this._sessionEstablished,
            sessionToken: this._sessionToken,
            statistics: this._statistics,
            passwordHash: UrlSafeBase64Encode(this._passwordHash),
            rootKey: UrlSafeBase64Encode(this._rootKey),
            privateKey: UrlSafeBase64Encode(this._privateKey),
            contactsKey: UrlSafeBase64Encode(this._contactsKey),
            sessionKey: UrlSafeBase64Encode(this._sessionKey)
        };

        return JSON.stringify(serializedData);
    }

    private _FindSessionToken(headers: Headers): string {
        const cookies = headers.getSetCookie();
        for (const cookie of cookies) {
            const [name, value] = cookie.split('=');
            if (name.trim() !== COOKIE_NAME) continue;
            const parts = value.split(';');
            let longest = '';
            for (const part of parts) if (part.length > longest.length) longest = part;
            return longest;
        }
        return '';
    }

    public async DecryptKeys(): Promise<void | CryptoGenericError> {

        // -- Key decryption
        try {
            const decompressedRootKey = Decompress(UrlSafeBase64Decode(this._encryptedSymetricRootKey));
            const decryptedRootKey = await Decrypt(decompressedRootKey, this._passwordHash);
            this._rootKey = UrlSafeBase64Decode(decryptedRootKey);

            const decompressedPrivateKey = Decompress(UrlSafeBase64Decode(this._encryptedAsymetricPrivateKey));
            const decryptedPrivateKey = await Decrypt(decompressedPrivateKey, this._rootKey);
            this._privateKey = UrlSafeBase64Decode(decryptedPrivateKey);

            const decompressedContactsKey = Decompress(UrlSafeBase64Decode(this._encryptedSymetricContactsKey));
            const decryptedContactsKey = await Decrypt(decompressedContactsKey, this._rootKey);
            this._contactsKey = UrlSafeBase64Decode(decryptedContactsKey);
        }

        catch (UnknownError) {
            return CryptoGenericError.fromUnknown(UnknownError, new CryptoGenericError(
                'Failed to decrypt keys'
            ));
        }

        // -- Integrity hash verification
        try {
            const integrityHash = await StandardIntegrityHash(this._rootKey, this._privateKey, this._contactsKey);
            if (integrityHash !== this._integrityHash) throw new CryptoGenericError('Integrity hash mismatch');
        }

        catch (UnknownError) {
            return CryptoGenericError.fromUnknown(UnknownError, new CryptoGenericError(
                'Failed to verify integrity hash'
            ));
        }

        this._passwordHash = new Uint8Array();
        this._sessionEstablished = true;
    }

    public async EncryptKey(key: Uint8Array | string): Promise<string> {
        if (typeof key !== 'string') key = UrlSafeBase64Encode(key);
        const encryptedKey = Compress(await Encrypt(key, this._rootKey));
        return UrlSafeBase64Encode(encryptedKey);
    }

    public async DecryptKey(key: Uint8Array | string): Promise<string> {
        if (typeof key !== 'string') key = UrlSafeBase64Encode(key);
        const decompressedKey = Decompress(UrlSafeBase64Decode(key));
        return await Decrypt(decompressedKey, this._rootKey);
    }

    public static ClearSessionCookie() {
        document.cookie = `${COOKIE_NAME}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
    }

    public get Statistics(): Statistics { return this._statistics; }
    public get SessionEstablished(): boolean { return this._sessionEstablished; }
    public get SessionKey(): Uint8Array { return this._sessionKey; }
    public get Token(): string { return this._sessionToken; }
    public get CookieToken(): string { return `${COOKIE_NAME}=${this._sessionToken}`; }
    public get IsTokenAuthenticated(): boolean { return this._sessionToken.length > 0; }
}

export default Session;

export {
    type SerializedSession
}
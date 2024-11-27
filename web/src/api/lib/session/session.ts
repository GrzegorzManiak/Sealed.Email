import {Statistics} from "./types";
import {ReturnedVerifyData} from "../auth/types";
import { Login } from "../auth/login";
import {Compress, Decompress, Decrypt, Encrypt} from "../symetric";
import {DecodeFromBase64} from "../common";
import {CryptoGenericError} from "../errors";
import {StandardIntegrityHash} from "../auth/register";
import {BigIntToByteArray, EncodeToBase64} from "gowl-client-lib";
import {COOKIE_NAME} from "../constants";

class Session {
    private readonly _encryptedSymetricRootKey: string;
    private readonly _encryptedAsymetricPrivateKey: string;
    private readonly _encryptedSymetricContactsKey: string;
    private readonly _integrityHash: string;

    private _sessionEstablished: boolean = false;
    private _sessionToken: string = '';

    private _passwordHash: Uint8Array;

    private _rootKey: Uint8Array;
    private _privateKey: Uint8Array;
    private _contactsKey: Uint8Array;
    private readonly _sessionKey: Uint8Array;

    private _statistics: Statistics;

    public constructor(awaitedLogin: Awaited<ReturnType<typeof Login>>, captureCookie: boolean = false) {
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
            const decompressedRootKey = Decompress(DecodeFromBase64(this._encryptedSymetricRootKey));
            const decryptedRootKey = await Decrypt(decompressedRootKey, this._passwordHash);
            this._rootKey = DecodeFromBase64(decryptedRootKey);

            const decompressedPrivateKey = Decompress(DecodeFromBase64(this._encryptedAsymetricPrivateKey));
            const decryptedPrivateKey = await Decrypt(decompressedPrivateKey, this._rootKey);
            this._privateKey = DecodeFromBase64(decryptedPrivateKey);

            const decompressedContactsKey = Decompress(DecodeFromBase64(this._encryptedSymetricContactsKey));
            const decryptedContactsKey = await Decrypt(decompressedContactsKey, this._rootKey);
            this._contactsKey = DecodeFromBase64(decryptedContactsKey);
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
        if (typeof key !== 'string') key = EncodeToBase64(key);
        const encryptedKey = Compress(await Encrypt(key, this._rootKey));
        return EncodeToBase64(encryptedKey);
    }

    public get Statistics(): Statistics { return this._statistics; }
    public get SessionEstablished(): boolean { return this._sessionEstablished; }
    public get SessionKey(): Uint8Array { return this._sessionKey; }
    public get Token(): string { return this._sessionToken; }
    public get CookieToken(): string { return `${COOKIE_NAME}=${this._sessionToken}`; }
    public get IsTokenAuthenticated(): boolean { return this._sessionToken.length > 0; }
}

export default Session;
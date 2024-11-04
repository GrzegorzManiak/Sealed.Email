import {Statistics} from "./types";
import {ReturnedVerifyData} from "../auth/types";
import { Login } from "../auth/login";
import {Decompress, Decrypt} from "../symetric";
import {DecodeFromBase64} from "../common";
import {CryptoGenericError} from "../errors";

class Session {
    private readonly _encryptedSymetricRootKey: string;
    private readonly _encryptedAsymetricPrivateKey: string;
    private readonly _encryptedSymetricContactsKey: string;

    private _passwordHash: Uint8Array;

    private _rootKey: Uint8Array;
    private _privateKey: Uint8Array;
    private _contactsKey: Uint8Array;

    private _statistics: Statistics;

    public constructor(awaitedLogin: Awaited<ReturnType<typeof Login>>) {
        if (awaitedLogin instanceof Error) throw awaitedLogin;
        const { verify, passwordHash } = awaitedLogin;
        this._passwordHash = passwordHash;

        this._encryptedSymetricRootKey = verify.encryptedSymmetricRootKey;
        this._encryptedAsymetricPrivateKey = verify.encryptedAsymmetricPrivateKey;
        this._encryptedSymetricContactsKey = verify.encryptedSymmetricContactsKey;

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

    public async DecryptKeys(): Promise<void | CryptoGenericError> {

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
            return CryptoGenericError.from_unknown(UnknownError, new CryptoGenericError(
                'Failed to decrypt keys'
            ));
        }

        this._passwordHash = new Uint8Array();
    }

    public get Statistics(): Statistics { return this._statistics; }
}

export default Session;
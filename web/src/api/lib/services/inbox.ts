import Session from "../session/session";
import {DomainFull} from "../api/domain";
import { DomainService } from "..";
import {InboxKeys} from "../api/inbox";
import {Decompress, Decrypt} from "../symetric";
import {DecodeFromBase64} from "../common";
import {GetCurve} from "gowl-client-lib";
import {CurrentCurve} from "../constants";

class InboxService {
    private readonly _inboxName: string;
    private _decryptedRootKey: Uint8Array;
    private _asymmetricPrivateKey: Uint8Array;
    private _asymmetricPublicKey: Uint8Array;


    public constructor(
        inboxName: string,
        decryptedRootKey: Uint8Array,
        asymmetricPrivateKey: Uint8Array,
        asymmetricPublicKey: Uint8Array
    ) {
        this._inboxName = inboxName;
        this._decryptedRootKey = decryptedRootKey;
        this._asymmetricPrivateKey = asymmetricPrivateKey;
        this._asymmetricPublicKey = asymmetricPublicKey;
    }

    public static async Decrypt(
        domainService: DomainService,
        inboxKeys: InboxKeys,
    ): Promise<InboxService> {
        const curve = GetCurve(CurrentCurve);
        const rootKey = DecodeFromBase64(await domainService.DecryptKey(inboxKeys.symmetricRootKey));

        // -- Decryption
        const decompressedPrivateKey = Decompress(DecodeFromBase64(inboxKeys.asymmetricPrivateKey));
        const asymmetricPrivateKey = DecodeFromBase64(await Decrypt(decompressedPrivateKey, rootKey));
        const asymmetricPublicKey = DecodeFromBase64(inboxKeys.asymmetricPublicKey);
        const calculatedPublicKey = curve.getPublicKey(asymmetricPrivateKey);
        const inboxName = await domainService.DecryptKey(inboxKeys.encryptedInboxName);

        // -- Sanity checks
        if (asymmetricPublicKey.toString() !== calculatedPublicKey.toString()) throw new Error("Invalid asymmetric public key");

        return new InboxService(inboxName, rootKey, asymmetricPrivateKey, asymmetricPublicKey);
    }

    public get InboxName(): string { return this._inboxName; }
}

export default InboxService;
import { GetCurve, Client } from 'gowl-client-lib';
import { CurrentCurve, Endpoints, ServerName } from '../constants';
import {Compress, Encrypt, NewKey} from '../symetric';
import { ProcessDetails, CalculateIntegrityHash, UrlSafeBase64Encode } from '../common';
import {ClientError, GenericError} from '../errors';
import { UserKeys } from './types';

async function StandardIntegrityHash(rootKey: Uint8Array, priv: Uint8Array, contactKey: Uint8Array): Promise<string> {
    return await CalculateIntegrityHash([rootKey, priv, contactKey]);
}

async function GenerateKeys(passwordHash: Uint8Array): Promise<UserKeys> {
    const rootKey = NewKey();
    const contactKey = NewKey();

    const curve = GetCurve(CurrentCurve);
    const priv = curve.utils.randomPrivateKey();
    const pub = curve.getPublicKey(priv);
    
    const EncryptedRootKey = await Encrypt(UrlSafeBase64Encode(rootKey), passwordHash);
    const EncryptedPrivateKey = await Encrypt(UrlSafeBase64Encode(priv), rootKey);
    const EncryptedContactKey = await Encrypt(UrlSafeBase64Encode(contactKey), rootKey);

    const integrityHash = await StandardIntegrityHash(rootKey, priv, contactKey);

    return {
        RootKey: rootKey,
        PublicKey: pub,
        PrivateKey: priv,
        ContactsKey: contactKey,
        EncryptedRootKey: Compress(EncryptedRootKey),
        EncryptedPrivateKey: Compress(EncryptedPrivateKey),
        EncryptedContactsKey: Compress(EncryptedContactKey),
        IntegrityHash: integrityHash
    }
}

async function SignUID(hashedUsername: string, keys: UserKeys): Promise<Uint8Array> {
    const curve = GetCurve(CurrentCurve);
    const bytes = new TextEncoder().encode(hashedUsername);

    try {
        const sig = curve.sign(bytes, keys.PrivateKey);
        return sig.toCompactRawBytes();
    } 

    catch (UnknownError) {
        throw new ClientError(
            'Failed to sign UID', 
            'Sorry, we were unable to sign your UID', 
            'SIGN-UID-FAIL'
        );
    }
}

async function UserCryptoSetup(username: string, password: string) {
    const { usernameHash, passwordHash } = await ProcessDetails(username, password);
    const keys = await GenerateKeys(passwordHash.hash);
    const proof = await SignUID(usernameHash, keys);
    return { keys, proof, usernameHash, passwordHash };
}

async function RegisterUser(username: string, password: string): Promise<UserKeys> {
    // -- Generate user keys
    const { keys, proof, usernameHash, passwordHash } = await UserCryptoSetup(username, password);

    // -- Register user
    const client = new Client(usernameHash, passwordHash.encoded, ServerName, CurrentCurve);
    const payload = await client.Register();
    if (payload instanceof Error) throw new ClientError(
        'Failed to register user',
        payload.message,
        'REG-OWL-FAIL'
    );

    // -- Send the request
    const response = await fetch(Endpoints.REGISTER[0], {
        method: Endpoints.REGISTER[1],
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            ...payload,
            TOS: true,
            Proof: UrlSafeBase64Encode(proof),
            PublicKey: UrlSafeBase64Encode(keys.PublicKey),
            EncryptedRootKey: UrlSafeBase64Encode(keys.EncryptedRootKey),
            EncryptedPrivateKey: UrlSafeBase64Encode(keys.EncryptedPrivateKey),
            EncryptedContactsKey: UrlSafeBase64Encode(keys.EncryptedContactsKey),
            IntegrityHash: keys.IntegrityHash
        })
    });

    if (!response.ok) throw GenericError.fromServerString(await response.text(), new ClientError(
        'Failed to register user',
        'Sorry, we were unable to register your account',
        'REG-REQ-FAIL'
    ));

    return keys;
}

export {
    RegisterUser,
    UserCryptoSetup,
    GenerateKeys,
    SignUID,
    StandardIntegrityHash
};
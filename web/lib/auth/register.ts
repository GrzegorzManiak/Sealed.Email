import { GetCurve, Hash, EncodeToBase64, Client, BigIntToByteArray } from 'gowl-client-lib';
import { CurrentCurve, Endpoints, ServerName } from '../constants';
import { Compress, Decompress, Encrypt } from '../symetric';
import { ProcessDetails, CalculateIntegrityHash } from '../common';
import { ClientError } from '../errors';
import { UserKeys } from './types';

async function StandardIntegrityHash(rootKey: Uint8Array, priv: Uint8Array, contactKey: Uint8Array): Promise<string> {
    return await CalculateIntegrityHash([rootKey, priv, contactKey]);
}

async function GenerateKeys(passwordHash: Uint8Array): Promise<UserKeys> {
    const curve = GetCurve(CurrentCurve);
    const rootKey = curve.utils.randomPrivateKey();
    const priv = curve.utils.randomPrivateKey();
    const contactKey = curve.utils.randomPrivateKey();
    const pub = curve.getPublicKey(priv);

    const EncryptedRootKey = await Encrypt(EncodeToBase64(rootKey), passwordHash);
    const EncryptedPrivateKey = await Encrypt(EncodeToBase64(priv), rootKey);
    const EncryptedContactKey = await Encrypt(EncodeToBase64(contactKey), rootKey);

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
};

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
};

async function UserCryptoSetup(username: string, password: string) {
    const { usernameHash, passwordHash } = await ProcessDetails(username, password);
    const keys = await GenerateKeys(passwordHash.hash);
    const proof = await SignUID(usernameHash, keys);
    return { keys, proof, usernameHash, passwordHash };
}

async function RegisterUser(username: string, password: string): Promise<ClientError | UserKeys> {

    try {
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
                Proof: EncodeToBase64(proof),
                PublicKey: EncodeToBase64(keys.PublicKey),
                EncryptedRootKey: EncodeToBase64(keys.EncryptedRootKey),
                EncryptedPrivateKey: EncodeToBase64(keys.EncryptedPrivateKey),
                EncryptedContactsKey: EncodeToBase64(keys.EncryptedContactsKey),
                IntegrityHash: keys.IntegrityHash
            })
        });

        if (!response.ok) {
            const errorText = await response.text();
            console.error('Registration failed:', errorText);
            throw new ClientError(
                'Failed to register user', 
                'Sorry, we were unable to register your account', 
                'REG-REQ-FAIL'
            );
        }

        return keys;
    }

    catch (UnknownError) {
        return ClientError.from_unknown(UnknownError, new ClientError(
            'Failed to register user', 
            'Sorry, we were unable to register your account', 
            'REG-CATCH'
        ));
    }
};

export {
    RegisterUser,
    UserCryptoSetup,
    GenerateKeys,
    SignUID,
    StandardIntegrityHash
};
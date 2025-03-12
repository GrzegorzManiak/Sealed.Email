import {BigIntToByteArray, EncodeToBase64, Hash} from "gowl-client-lib";
import {ServerName} from "./constants";

function UrlSafeBase64Encode(data: Uint8Array | bigint): string {
    if (typeof data === 'bigint') data = BigIntToByteArray(data);
    return EncodeToBase64(data).replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '');
}

function UrlSafeBase64Decode(data: string): Uint8Array {
    const base64 = data.replace(/-/g, '+').replace(/_/g, '/');
    const padded = base64.padEnd(base64.length + (4 - (base64.length % 4)) % 4, '='); // -- Padding is needed for JS Built-in atob
    const binaryString = atob(padded);
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) bytes[i] = binaryString.charCodeAt(i);
    return bytes;
}

async function Argon2Hash(username: string, password: string) {

    // -- TODO -- FIND A LIBRARY THAT SUPPORTS ARGON2ID
    // import argon2 from 'argon2';
    // return await argon2.hash({
    //     pass: password,
    //     salt: username,
    //     time: 5,
    //     mem: 1024 * 128,
    //     hashLen: 32
    // });
    // WILL HAVE TO USE DIFFERENT HASH FOR NOW
    const hashed = BigIntToByteArray(await Hash(password + username));
    return {
        hash: hashed,
        encoded: UrlSafeBase64Encode(hashed)
    }
}

async function CalculateIntegrityHash(keys: Array<Uint8Array>): Promise<string> {
    const joinedKeys = new Uint8Array(keys.reduce((acc, key) => acc + key.length, 0));
    return UrlSafeBase64Encode(await Hash(UrlSafeBase64Encode(joinedKeys)));
}

async function ProcessDetails(username: string, password: string) {
    const usernameHash: string = UrlSafeBase64Encode(BigIntToByteArray(await Hash(username + ServerName)));
    const passwordHash = await Argon2Hash(usernameHash, password);
    return { usernameHash, passwordHash };
}

function SplitEmail(email: string): { username: string, domain: string } {
    email = email.trim().toLowerCase();
    const parts = email.split('@', 2);
    return { username: parts[0], domain: parts[1] };
}

function EnsureFqdn(domain: string): string {
    domain = domain.trim().toLowerCase();
    if (domain.endsWith('.')) return domain;
    else return domain + '.';
}

async function HashInboxEmail(email: string): Promise<string> {
    const { domain, username } = SplitEmail(email);
    if (!domain || !username) throw new Error("Invalid email");
    const userHash = await Hash(`${username}@${EnsureFqdn(domain)}`);
    return `${UrlSafeBase64Encode(userHash)}@${domain}`;
}

export {
    Argon2Hash,
    ProcessDetails,
    UrlSafeBase64Encode,
    UrlSafeBase64Decode,
    CalculateIntegrityHash,
    HashInboxEmail,
    ServerName,
    SplitEmail,
    EnsureFqdn
};
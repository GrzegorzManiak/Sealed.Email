import {ALG, IV_LENGTH} from "./constants";

async function Encrypt(text: string, key: Uint8Array): Promise<{ iv: number[], data: number[] }> {
    const encoder = new TextEncoder();
    const data = encoder.encode(text);

    if (key.length !== 32) throw new Error("Key must be 32 bytes");

    const iv = crypto.getRandomValues(new Uint8Array(IV_LENGTH));

    const algorithm = { name: ALG, iv };
    const cryptoKey = await crypto.subtle.importKey("raw", key, algorithm, false, ["encrypt"]);
    const encryptedData = await crypto.subtle.encrypt(algorithm, cryptoKey, data);

    return { iv: Array.from(iv), data: Array.from(new Uint8Array(encryptedData)) };
}


async function Decrypt(encryptedData: { iv: number[], data: number[] }, key: Uint8Array): Promise<string> {
    const { iv, data } = encryptedData;
    const ivArray = new Uint8Array(iv);
    const dataArray = new Uint8Array(data);

    if (key.length !== 32) throw new Error("Key must be 32 bytes");

    const algorithm = { name: ALG, iv: ivArray };
    const cryptoKey = await crypto.subtle.importKey('raw', key, algorithm, false, ['decrypt']);
    const decrypted = await crypto.subtle.decrypt(algorithm, cryptoKey, dataArray);

    const decoder = new TextDecoder();
    return decoder.decode(decrypted);
}

function Compress(encryptedData: { iv: number[], data: number[] }): Uint8Array {
    const { iv, data } = encryptedData;
    const ivArray = new Uint8Array(iv);
    const dataArray = new Uint8Array(data);
    return new Uint8Array([...ivArray, ...dataArray]);
}

function Decompress(compressedData: Uint8Array): { iv: number[], data: number[] } {
    const iv = compressedData.slice(0, IV_LENGTH);
    const data = compressedData.slice(IV_LENGTH);
    return { iv: Array.from(iv), data: Array.from(data) };
}

function NewKey(length = 32): Uint8Array {
    return crypto.getRandomValues(new Uint8Array(length));
}
  
export {
    Encrypt,
    Decrypt,
    Compress,
    Decompress,
    NewKey
}
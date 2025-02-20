import {ALG, IVLength, DefaultKeyLength} from "./constants";

async function Encrypt(text: string, key: Uint8Array): Promise<{ iv: number[], data: number[] }> {
    const encoder = new TextEncoder();
    const data = encoder.encode(text);

    if (key.length !== DefaultKeyLength) throw new Error(`Key must be ${DefaultKeyLength} bytes, got ${key.length}`);

    const iv = crypto.getRandomValues(new Uint8Array(IVLength));

    const algorithm = { name: ALG, iv };
    const cryptoKey = await crypto.subtle.importKey("raw", key, algorithm, false, ["encrypt"]);
    const encryptedData = await crypto.subtle.encrypt(algorithm, cryptoKey, data);

    return { iv: Array.from(iv), data: Array.from(new Uint8Array(encryptedData)) };
}

async function Decrypt(encryptedData: { iv: number[], data: number[] }, key: Uint8Array): Promise<string> {
    const { iv, data } = encryptedData;
    const ivArray = new Uint8Array(iv);
    const dataArray = new Uint8Array(data);

    if (key.length !== DefaultKeyLength) throw new Error(`Key must be ${DefaultKeyLength} bytes, got ${key.length}`);

    const algorithm = { name: ALG, iv: ivArray };
    const cryptoKey = await crypto.subtle.importKey('raw', key, algorithm, false, ['decrypt']);
    const decrypted = await crypto.subtle.decrypt(algorithm, cryptoKey, dataArray);

    const decoder = new TextDecoder();
    return decoder.decode(decrypted);
}

function Compress(iv: Uint8Array, data: Uint8Array): Uint8Array {
    const ivLen = new Uint8Array([iv.length]);
    return new Uint8Array([...ivLen, ...iv, ...data]);
}

function Decompress(compressedData: Uint8Array): { iv: Uint8Array; data: Uint8Array } {
    if (compressedData.length < 2) {
        throw new Error("Invalid compressed data: too short");
    }

    const ivLen = compressedData[0];
    if (compressedData.length < 1 + ivLen) {
        throw new Error("Invalid compressed data: truncated IV");
    }

    const iv = compressedData.slice(1, 1 + ivLen);
    const data = compressedData.slice(1 + ivLen);
    return { iv, data };
}

function NewKey(length = DefaultKeyLength): Uint8Array {
    return crypto.getRandomValues(new Uint8Array(length));
}
  
export {
    Encrypt,
    Decrypt,
    Compress,
    Decompress,
    NewKey
}
async function Encrypt(text: string, key: Uint8Array): Promise<{ iv: number[], data: number[] }> {
    const encoder = new TextEncoder();
    const data = encoder.encode(text);
    const iv = crypto.getRandomValues(new Uint8Array(16)); 
  
    const cryptoKey = await crypto.subtle.importKey(
        'raw', key, { name: 'AES-GCM' }, false, ['encrypt']);
  
    const encrypted = await crypto.subtle.encrypt(
        { name: 'AES-GCM', iv: iv }, cryptoKey, data);
  
    return { iv: Array.from(iv), data: Array.from(new Uint8Array(encrypted)) };
};
  
async function Decrypt(encryptedData: { iv: number[], data: number[] }, key: Uint8Array): Promise<string> {
    const { iv, data } = encryptedData;
    const ivArray = new Uint8Array(iv);
    const dataArray = new Uint8Array(data);
  
    const cryptoKey = await crypto.subtle.importKey(
        'raw', key, { name: 'AES-GCM' }, false, ['decrypt']);
  
    const decrypted = await crypto.subtle.decrypt(
      { name: 'AES-GCM', iv: ivArray }, cryptoKey, dataArray);
  
    const decoder = new TextDecoder();
    return decoder.decode(decrypted);
};

function Compress(encryptedData: { iv: number[], data: number[] }): Uint8Array {
    const { iv, data } = encryptedData;
    const ivArray = new Uint8Array(iv);
    const dataArray = new Uint8Array(data);
    return new Uint8Array([ ...ivArray, ...dataArray ]);
};

function Decompress(compressedData: Uint8Array): { iv: number[], data: number[] } {
    const iv = compressedData.slice(0, 16);
    const data = compressedData.slice(16);
    return { iv: Array.from(iv), data: Array.from(data) };
}
  
export {
    Encrypt,
    Decrypt,
    Compress,
    Decompress
}
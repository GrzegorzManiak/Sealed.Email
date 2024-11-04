import argon2 from 'argon2-browser';
import {BigIntToByteArray, EncodeToBase64, Hash} from "gowl-client-lib";
import {ServerName} from "./constants";

async function Argon2Hash(username: string, password: string) {
    return await argon2.hash({
        pass: password,
        salt: username,
        time: 5,
        mem: 1024 * 128,
        hashLen: 32
    });
}

async function ProcessDetails(username: string, password: string) {
    const usernameHash: string = EncodeToBase64(BigIntToByteArray(await Hash(username + ServerName)));
    const passwordHash = await Argon2Hash(usernameHash, password);
    return { usernameHash, passwordHash };
}

function DecodeFromBase64(base64String: string): Uint8Array {
    const binaryString = atob(base64String);
    const len = binaryString.length;
    const bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) bytes[i] = binaryString.charCodeAt(i);
    return bytes;
}

export {
    Argon2Hash,
    ProcessDetails,
    DecodeFromBase64,
    ServerName
};
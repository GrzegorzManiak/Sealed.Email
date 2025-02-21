import {CurrentCurve} from "./constants";
import {BigIntToByteArray, BytesToBigInt, GetCurve, Hash} from "gowl-client-lib";
import {UrlSafeBase64Decode, UrlSafeBase64Encode} from "./common";
import * as Sym from "./symetric";

function GenerateKeyPair(curve= GetCurve(CurrentCurve)) {
	const priv = curve.utils.randomPrivateKey();
	const pub = curve.getPublicKey(priv);
	return {priv, pub};
}

async function SignData(data: string | Uint8Array, priv: Uint8Array, curve= GetCurve(CurrentCurve)) {
	const encodedData = typeof data === 'string' ? new TextEncoder().encode(data) : data;
	const signature = curve.sign(encodedData, priv);
	return UrlSafeBase64Encode(signature.toCompactRawBytes());
}

async function VerifySignature(data: string | Uint8Array, signature: string, pub: Uint8Array, curve = GetCurve(CurrentCurve)) {
	const encodedData = typeof data === 'string' ? new TextEncoder().encode(data) : data;
	return curve.verify(UrlSafeBase64Decode(signature), encodedData, pub);
}

async function SharedKey(priv: Uint8Array, pub: Uint8Array, curve = GetCurve(CurrentCurve)) {
	return curve.getSharedSecret(BytesToBigInt(priv), pub);
}

async function Encrypt(
	data: Uint8Array,
	sharedKey: Uint8Array
): Promise<string> {
	const stringData = UrlSafeBase64Encode(data);
	const sizedKey = await Hash(sharedKey);
	const encrypted = await Sym.Encrypt(stringData, BigIntToByteArray(sizedKey));
	const bytes = Sym.Compress(encrypted);
	return UrlSafeBase64Encode(bytes);
}

async function Decrypt(
	cipherText: Uint8Array,
	sharedKey: Uint8Array
): Promise<string> {
	const decompressed = Sym.Decompress(cipherText);
	const sizedKey = await Hash(sharedKey);
	const decrypted = await Sym.Decrypt(decompressed, BigIntToByteArray(sizedKey));
	return Buffer.from(decrypted, 'base64').toString();
}

export {
	GenerateKeyPair,
	SignData,
	VerifySignature,
	SharedKey,
	Encrypt,
	Decrypt,
}
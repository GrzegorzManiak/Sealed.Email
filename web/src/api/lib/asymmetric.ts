import {CurrentCurve} from "./constants";
import {BigIntToByteArray, BytesToBigInt, GetCurve, Hash} from "gowl-client-lib";
import {UrlSafeBase64Decode} from "./common";
import * as Sym from "./symetric";

function GenerateKeyPair(curve= GetCurve(CurrentCurve)) {
	const priv = curve.utils.randomPrivateKey();
	const pub = curve.getPublicKey(priv);
	return {priv, pub};
}

function UrlSafeBase64Encode(data: Uint8Array | string) {
	if (typeof data === 'string') return Buffer.from(data).toString('base64');
	else return Buffer.from(data).toString('base64');
}

async function SignData(data: string | Uint8Array, priv: Uint8Array, curve= GetCurve(CurrentCurve)) {
	const encodedData = data instanceof Uint8Array ? data : new TextEncoder().encode(data);
	const signature = curve.sign(encodedData, priv);
	const bytes = signature.toCompactRawBytes();
	return UrlSafeBase64Encode(bytes);
}

async function VerifySignature(data: string | Uint8Array, signature: string, pub: Uint8Array, curve = GetCurve(CurrentCurve)) {
	const encodedData = data instanceof Uint8Array ? data : new TextEncoder().encode(data);
	return curve.verify(UrlSafeBase64Decode(signature), encodedData, pub);
}

async function SharedKey(priv: Uint8Array, pub: Uint8Array, curve = GetCurve(CurrentCurve)) {
	return curve.getSharedSecret(BytesToBigInt(priv), pub);
}

async function Encrypt(
	data: string | Uint8Array,
	sharedKey: Uint8Array
): Promise<string> {
	const stringData = UrlSafeBase64Encode(data);
	const sizedKey = await Hash(sharedKey);
	const encrypted = await Sym.Encrypt(stringData, BigIntToByteArray(sizedKey));
	const bytes = Sym.Compress(encrypted);
	return UrlSafeBase64Encode(bytes);
}

async function Decrypt(
	data: string | Uint8Array,
	sharedKey: Uint8Array
): Promise<string> {
	const bytes = data instanceof Uint8Array ? data : UrlSafeBase64Decode(data);
	const decompressed = Sym.Decompress(bytes);
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
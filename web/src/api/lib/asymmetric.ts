import {CurrentCurve} from "./constants";
import {NewKey} from "./symetric";
import {BigIntToByteArray, BytesToBigInt, GetCurve, Hash} from "gowl-client-lib";
import {DecodeFromBase64} from "./common";
import * as Sym from "./symetric";

function GenerateKeyPair(curve= GetCurve(CurrentCurve)) {
	const priv = curve.utils.randomPrivateKey();
	const pub = curve.getPublicKey(priv);
	return {priv, pub};
}

function EncodeToBase64(data: Uint8Array | string) {
	if (typeof data === 'string') return Buffer.from(data).toString('base64');
	else return Buffer.from(data).toString('base64');
}

async function SignData(data: string | Uint8Array, priv: Uint8Array, curve= GetCurve(CurrentCurve)) {
	const encodedData = data instanceof Uint8Array ? data : new TextEncoder().encode(data);
	const signature = curve.sign(encodedData, priv);
	const bytes = signature.toCompactRawBytes();
	return EncodeToBase64(bytes);
}

async function VerifySignature(data: string | Uint8Array, signature: string, pub: Uint8Array, curve = GetCurve(CurrentCurve)) {
	const encodedData = data instanceof Uint8Array ? data : new TextEncoder().encode(data);
	return curve.verify(DecodeFromBase64(signature), encodedData, pub);
}

async function SharedKey(priv: Uint8Array, pub: Uint8Array, curve = GetCurve(CurrentCurve)) {
	return curve.getSharedSecret(BytesToBigInt(priv), pub);
}

async function Encrypt(
	data: string | Uint8Array,
	sharedKey: Uint8Array
): Promise<string> {
	const stringData = EncodeToBase64(data);
	const sizedKey = await Hash(sharedKey);
	const encrypted = await Sym.Encrypt(stringData, BigIntToByteArray(sizedKey));
	const bytes = Sym.Compress(encrypted);
	return EncodeToBase64(bytes);
}

async function Decrypt(
	data: string | Uint8Array,
	sharedKey: Uint8Array
): Promise<string> {
	const bytes = data instanceof Uint8Array ? data : DecodeFromBase64(data);
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
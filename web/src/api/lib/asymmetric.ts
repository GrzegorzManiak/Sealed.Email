import {CurrentCurve, IVLength} from "./constants";
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

export function normalizeKey(publicKey: Uint8Array): Uint8Array {
	// -- Already uncompressed
	if (publicKey.length === 65 && publicKey[0] === 0x04) return publicKey;

	// -- Compressed key, decompress it (With the sign byte 0x02 or 0x03 / without sign byte)
	else if ((publicKey.length === 33 && (publicKey[0] === 0x02 || publicKey[0] === 0x03)) || publicKey.length === 32) {
		const curve = GetCurve(CurrentCurve);
		const point = curve.ProjectivePoint.fromHex(publicKey);
		return point.toRawBytes(false);
	}

	throw new Error(`Invalid public key length: ${publicKey.length}`);
}


function CompressKey(publicKey: Uint8Array): Uint8Array {
	const curve = GetCurve(CurrentCurve);
	const point = curve.ProjectivePoint.fromHex(publicKey);
	return point.toRawBytes(true);
}

async function Encrypt(
	data: Uint8Array,
	sharedKey: Uint8Array
): Promise<string> {
	const stringData = UrlSafeBase64Encode(data);
	const encrypted = await Sym.Encrypt(stringData, sharedKey);
	const bytes = Sym.Compress(encrypted);
	return UrlSafeBase64Encode(bytes);
}

async function Decrypt(
	cipherText: Uint8Array,
	ephemeralKey: Uint8Array,
	privateKey: Uint8Array
): Promise<string> {
	ephemeralKey = CompressKey(ephemeralKey);

	console.log("ephemeralKey", ephemeralKey);
	console.log("privateKey", privateKey);

	const curve = GetCurve(CurrentCurve);
	const sharedSecret = curve.getSharedSecret(BytesToBigInt(privateKey), ephemeralKey, true);
	const hashedSecret = BigIntToByteArray(await Hash(sharedSecret));
	console.log("hashedSecret", hashedSecret);
	const decompressed = Sym.Decompress(cipherText);
	return await Sym.Decrypt(decompressed, hashedSecret);
}

function Decompress(compressedData: Uint8Array): { ephemeralKey: Uint8Array, cipherText: Uint8Array } {
	const keySize = 1 + 2 * (GetCurve(CurrentCurve).CURVE.nBitLength / 8);
	if (compressedData.length < keySize + IVLength) throw new Error("Invalid encrypted data format");
	const ephemeralKey = compressedData.slice(0, keySize);
	const cipherText = compressedData.slice(keySize);
	return { ephemeralKey, cipherText };
}

export {
	GenerateKeyPair,
	Decompress,
	SignData,
	VerifySignature,
	SharedKey,
	Encrypt,
	Decrypt,
}
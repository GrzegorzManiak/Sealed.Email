import {CurrentCurve, IVLength} from "./constants";
import {BigIntToByteArray, BytesToBigInt, GetCurve, Hash} from "gowl-client-lib";
import {UrlSafeBase64Decode, UrlSafeBase64Encode} from "./common";
import * as Sym from "./symetric";
import {sha256} from "js-sha256";

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
	publicKey: Uint8Array,
	privateKey: Uint8Array,
	ciphertext: Uint8Array,
	curve = GetCurve(CurrentCurve)
): Promise<string> {
	const pubPoint = curve.ProjectivePoint.fromHex(publicKey);
	const sharedSecretPoint = pubPoint.multiply(BytesToBigInt(privateKey));
	const sharedX = sharedSecretPoint.toRawBytes(true).slice(1);
	const sharedKey = BigIntToByteArray(await Hash(BytesToBigInt(sharedX)));
	return await Sym.Decrypt(Sym.Decompress(ciphertext), sharedKey);
}

function ExtractEphemeralPubKey(
	encryptedData: Uint8Array
): { ephemeralPub: Uint8Array, keyLength: number, ciphertext: Uint8Array } {
	if (encryptedData.length < 2 + IVLength) {
		throw new Error("Invalid encrypted data format");
	}

	// Extract key length from first 2 bytes
	const keyLength = (encryptedData[0] << 8) | encryptedData[1];

	if (encryptedData.length < 2 + keyLength + IVLength) {
		throw new Error("Invalid encrypted data format");
	}

	// Extract ephemeral key bytes using the length we just read
	const keyBytes = encryptedData.slice(2, 2 + keyLength);

	// Return the ephemeral key, its length, and the remaining ciphertext
	return {
		ephemeralPub: keyBytes,
		keyLength: keyLength,
		ciphertext: encryptedData.slice(2 + keyLength)
	};
}

export {
	GenerateKeyPair,
	SignData,
	ExtractEphemeralPubKey,
	Encrypt,
	Decrypt
}
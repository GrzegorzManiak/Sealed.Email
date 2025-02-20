import { SupportedCurves } from "gowl-client-lib";
import { secp256k1 } from "@noble/curves/secp256k1";

const ROOT = 'http://localhost:2095/api';
const COOKIE_NAME = 'NES-DEV';

const Endpoints: { [key: string]: [string, string] } = {
    REGISTER:       [`${ROOT}/register`, 'POST'],
    LOGIN_INIT:     [`${ROOT}/login/init`, 'PUT'],
    LOGIN_VERIFY:   [`${ROOT}/login/verify`, 'PUT'],

    DOMAIN_ADD:     [`${ROOT}/domain/add`, 'POST'],
    DOMAIN_GET:     [`${ROOT}/domain/get`, 'GET'],
    DOMAIN_DELETE:  [`${ROOT}/domain/delete`, 'DELETE'],
    DOMAIN_LIST:    [`${ROOT}/domain/list`, 'GET'],
    DOMAIN_MODIFY:  [`${ROOT}/domain/modify`, 'PUT'],
    DOMAIN_REFRESH:  [`${ROOT}/domain/refresh`, 'PUT'],

    EMAIL_SEND_PLAIN: [`${ROOT}/email/send/plain`, 'POST'],
    EMAIL_SEND_ENCRYPTED: [`${ROOT}/email/send/encrypted`, 'POST'],
    EMAIL_LIST: [`${ROOT}/email/list`, 'GET'],
    EMAIL_GET: [`${ROOT}/email/get`, 'GET'],
    EMAIL_DATA: [`${ROOT}/email/data`, 'GET'],
};

const ServerName: string = 'NoiseEmailServer>V1.0.0';
const CurrentCurve = SupportedCurves.P256;
const ALG = 'AES-GCM'
const IVLength = 12;
const DefaultKeyLength = 32;

export {
    CurrentCurve,
    ServerName,
    Endpoints,
    COOKIE_NAME,
    DefaultKeyLength,
    IVLength,
    ROOT,
    ALG
}
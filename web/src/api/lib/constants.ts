import { SupportedCurves } from "gowl-client-lib";

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
};

const ServerName: string = 'NoiseEmailServer>V1.0.0';
const CurrentCurve: SupportedCurves = SupportedCurves.P256;

export {
    CurrentCurve,
    ServerName,
    Endpoints,
    COOKIE_NAME,
    ROOT
}
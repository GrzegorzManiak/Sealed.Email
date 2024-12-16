import { SupportedCurves } from "gowl-client-lib";

const ROOT = 'http://localhost:2095/api';
const COOKIE_NAME = 'NES-DEV';

const Endpoints: { [key: string]: [string, string] } = {
    REGISTER:       [`${ROOT}/register`, 'POST'],
    LOGIN_INIT:     [`${ROOT}/login/init`, 'PUT'],
    LOGIN_VERIFY:   [`${ROOT}/login/verify`, 'PUT'],

    DOMAIN_ADD:     [`${ROOT}/domain/add`, 'POST'],
    DOMAIN_GET:     [`${ROOT}/domain/get`, 'PUT'],
    DOMAIN_DELETE:  [`${ROOT}/domain/delete`, 'DELETE'],
    DOMAIN_LIST:    [`${ROOT}/domain/list`, 'PUT'],
    DOMAIN_MODIFY:  [`${ROOT}/domain/modify`, 'PUT'],
    DOMAIN_REFRESH:  [`${ROOT}/domain/refresh`, 'PUT'],

    INBOX_ADD:      [`${ROOT}/inbox/add`, 'POST'],
    INBOX_DELETE:   [`${ROOT}/inbox/delete`, 'DELETE'],
    INBOX_LIST:     [`${ROOT}/inbox/list`, 'GET'],
    INBOX_MODIFY:   [`${ROOT}/inbox/modify`, 'PUT'],
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
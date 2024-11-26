import { SupportedCurves } from "gowl-client-lib";

const ROOT = 'http://localhost:2095/api';
const COOKIE_NAME = 'NES-DEV';

const Endpoints = {
    REGISTER:       [`${ROOT}/register`, 'POST'],
    LOGIN_INIT:     [`${ROOT}/login/init`, 'PUT'],
    LOGIN_VERIFY:   [`${ROOT}/login/verify`, 'PUT'],

    DOMAIN_ADD:     [`${ROOT}/domain/add`, 'POST'],
    DOMAIN_DELETE:  [`${ROOT}/domain/delete`, 'DELETE'],
    DOMAIN_LIST:    [`${ROOT}/domain/list`, 'PUT'],
    DOMAIN_MODIFY:  [`${ROOT}/domain/modify`, 'PUT'],
    DOMAIN_REFRESH:  [`${ROOT}/domain/refresh`, 'PUT'],
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
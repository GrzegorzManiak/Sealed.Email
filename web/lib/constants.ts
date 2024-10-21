import { SupportedCurves } from "gowl-client-lib";

const ROOT = 'http://localhost:8080';

const Endpoints = {
    REGISTER:       [`${ROOT}/register`, 'POST'],
    LOGIN_INIT:     [`${ROOT}/login/init`, 'PUT'],
    LOGIN_VERIFY:   [`${ROOT}/login/verify`, 'PUT']
};

const ServerName: string = 'NoiseEmailServer>V1.0.0';
const CurrentCurve: SupportedCurves = SupportedCurves.P256;

export {
    CurrentCurve,
    ServerName,
    Endpoints,
    ROOT
}
import { v4 as uuidv4 } from 'uuid';
import {throwToast} from "../../lib/toasts";

class GenericError extends Error {
    protected readonly _id: string = uuidv4();
    protected readonly _message: string;
    protected readonly _title: string = this.constructor.name;

    public constructor(
        message: string,
        title: string
    ) {
        super(message);
        this._message = message;
        if (title) this._title = title;
    }



    public toast() {
        throwToast(this._title, this._message);
    }



    public get id(): string { return this._id; }
    public get message(): string { return this._message; }
    public get title(): string { return this._title; }



    public static fromError (error: Error): GenericError {
        return new GenericError(error.message, error.name);
    };

    public static fromUnknown (
        error: unknown,
        elseError: GenericError = new GenericError("Sorry, we were unable to process your request", "Oops! An unknown error occurred"),
    ): GenericError {
        if (error instanceof GenericError) return error;
        else if (error instanceof Error) return GenericError.fromError(error);
        else if (typeof error === 'string') return new GenericError(error, "Oops! An unknown error occurred");
        return elseError;
    };

    public static fromServer(
        obj: { msg: string; title: string; code: number; error: boolean },
        elseError: GenericError = new GenericError("Sorry, we were unable to process your request", "Oops! An unknown error occurred"),
    ): GenericError {
        if (typeof obj !== 'object' || obj === null) return elseError;
        const requiredProps = ['msg', 'title', 'fault', 'error'];
        if (requiredProps.every((key) => key in obj)) return new GenericError(obj.msg, obj.title);
        return elseError;
    };

    public static fromServerString (
        str: string,
        elseError: GenericError = new GenericError("Sorry, we were unable to process your request", "Oops! An unknown error occurred"),
    ): GenericError {
        try {
            const obj = JSON.parse(str);
            return GenericError.fromServer(obj, elseError);
        } catch {
            return elseError;
        }
    }
}
class ClientError extends GenericError {
    public constructor(title: string, message: string, code: string) {
        super(message, title);
    }
}

class CryptoGenericError extends GenericError {
    public constructor(message: string) {
        super(message, 'Failure in cryptographic operation');
    }
}

export {
    GenericError,
    CryptoGenericError,
    ClientError
}

export default GenericError;
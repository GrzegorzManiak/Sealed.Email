import { v4 as uuidv4 } from 'uuid';

type SerializedGenericError = {
    id: string,
    message: string,
    hint?: string,
    code: number,
    data: object,
    title: string,
    errors: Array<SerializedGenericError>
}

class GenericError extends Error {
    private _other_errors: Map<string, GenericError> = new Map();

    protected readonly _id: string = uuidv4();
    protected readonly _message: string;
    protected readonly _code: number;
    protected readonly _title: string = this.constructor.name;
    protected _data: object = {};
    protected _hint: string = '';

    public constructor(
        message: string,
        code: number,
        title: string
    ) {
        super(message);
        this._message = message;
        this._code = code;
        if (title) this._title = title;
    }

    // -- Function that should be overridden by child classes
    public serialize = (): SerializedGenericError => {
        return {
            id: this._id,
            message: this._message,
            code: this._code,
            data: this._data,
            hint: this._hint,
            title: this._title,
            errors: this.errors.map((error) => error.serialize())
        };
    };  



    public static deserialize = (serialized: SerializedGenericError): GenericError => {
        return new DeserializedGenericError(serialized);
    }
        


    public get errors(): Array<GenericError> {
        const existing_ids = new Set<string>();
        return Array.from(this._other_errors.values()).filter((error) => {
            if (existing_ids.has(error.id)) return false;
            existing_ids.add(error.id);
            return true;
        });
    }


    public toString(): string {
        return JSON.stringify(this.serialize());
    }

    public set data(data: object) { this._data = data; }
    public get data(): object { return this._data; }

    public get hint(): string { return this._hint; }
    public set hint(hint: string) { this._hint = hint; }


    public get id(): string { return this._id; }
    public get message(): string { return this._message; }
    public get code(): number { return this._code; }
    public get title(): string { return this._title; }


    public add_error = (error: GenericError): void => {
        // -- Ensure that the error is not already in the map
        if (this._other_errors.has(error.id)) return;
        
        // -- Ensure that we are not adding ourselves to the map
        if (error.id === this._id) return;

        // -- Add the error to the map
        this._other_errors.set(error.id, error);
    };



    public static is_error = (error: unknown): boolean => {
        // -- If the error is anything bar a class its 100% not a generic error
        if (typeof error !== 'object' || error === null) return false;
        return (
            error instanceof GenericError || 
            error instanceof Error ||
            ['id', 'message', 'code', 'title', 'data'].every((key) => key in error)
        );
    };


    
    public static from_error = (error: Error): GenericError => {
        return new GenericError(error.message, 500, 'Error');
    };



    public static from_unknown = (
        error: unknown,
        else_error: GenericError = new GenericError("Sorry, we were unable to process your request", 500, "Oops! An unknown error occurred"),
        hint?: string
    ): GenericError => {
        let return_error = else_error;
        
        if (error instanceof GenericError) return_error = error;
        else if (error instanceof Error) return_error = GenericError.from_error(error);
        else if (typeof error === 'string') return_error = new GenericError(error, 500, 'Oops! An unknown error occurred');
        if (hint) return_error.hint = hint;

        return return_error;
    };

    public static from_server = (
        obj: { msg: string, title: string, code: number, error: boolean },
        backup_error: GenericError = new GenericError("Sorry, we were unable to process your request", 500, "Oops! An unknown error occurred"),
    ): GenericError => {
        if (typeof obj !== 'object' || obj === null) return backup_error;
        const requiredProps = ['msg', 'title', 'fault', 'error'];
        if (requiredProps.every((key) => key in obj)) return new GenericError(obj.msg, 400, obj.title);
        return backup_error;
    };

    public static from_server_string = (
        str: string,
        backup_error: GenericError = new GenericError("Sorry, we were unable to process your request", 500, "Oops! An unknown error occurred"),
    ): GenericError => {
        try {
            const obj = JSON.parse(str);
            return GenericError.from_server(obj, backup_error);
        } catch {
            return backup_error;
        }
    }
}

class DeserializedGenericError extends GenericError {
    protected readonly _id: string;
    protected readonly _message: string;
    protected readonly _code: number;
    protected readonly _title: string;
    protected _data: object;
    protected _hint: string;

    public constructor(serialized: SerializedGenericError) {
        super(serialized.message, serialized.code, serialized.title);
        this._id = serialized.id;
        this._message = serialized.message;
        this._code = serialized.code;
        this._data = serialized.data;
        this._hint = serialized.hint || '';
        this._title = serialized.title;
        serialized.errors.forEach((error) => 
            this.add_error(new DeserializedGenericError(error)));
    }
}

class ClientError extends GenericError {
    public constructor(title: string, message: string, code: string) {
        super(message, 400, title);
        this._hint = code;
    }
}

class CryptoGenericError extends GenericError {
    public constructor(message: string) {
        super(message, 400, 'Failure in cryptographic operation');
    }
}

export {
    SerializedGenericError,
    GenericError,
    DeserializedGenericError,
    CryptoGenericError,
    ClientError
}
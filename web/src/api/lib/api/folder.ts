import { parseDomain, ParseResultType } from "parse-domain";
import Session from "../session/session";
import {Endpoints} from "../constants";
import {ClientError} from "../errors";
import {NewKey} from "../symetric";
import {HandleRequest} from "./common";
import * as Asym from "../asymmetric";
import * as Sym from "../symetric";
import { UrlSafeBase64Encode } from "../common";
import {DomainFull, DomainRefID} from "./domain";
import {DomainService} from "../index";

type Folder = {
	encryptedFolderName: string;
	folderID: string;
}

type FolderListResponse = {
	folders: Folder[];
	total: number;
}

const FolderList = async (session: Session, domainID: DomainRefID): Promise<FolderListResponse> => HandleRequest<FolderListResponse>({
	session,
	query: { domainID },
	endpoint: Endpoints.FOLDER_LIST,
	fallbackError: new ClientError(
		'Failed to get folders',
		'Sorry, we were unable to get your folders! Please try again later.',
		'FOLDER-LIST-FAIL'
	),
});

const FolderCreate = async (session: Session, domain: DomainService, folderName: string): Promise<Folder> => {
	const encryptedFolderName = await domain.EncryptData(folderName);
	return HandleRequest<Folder>({
		session,
		body: { encryptedFolderName, domainID: domain.DomainID },
		endpoint: Endpoints.FOLDER_CREATE,
		fallbackError: new ClientError(
			'Failed to create folder',
			'Sorry, we were unable to create your folder! Please try again later.',
			'FOLDER-CREATE-FAIL'
		),
	});
}

const FolderDelete = async (session: Session, domainID: DomainRefID, folderID: string): Promise<void> => HandleRequest<void>({
	session,
	query: { domainID, folderID },
	endpoint: Endpoints.FOLDER_DELETE,
	fallbackError: new ClientError(
		'Failed to delete folder',
		'Sorry, we were unable to delete your folder! Please try again later.',
		'FOLDER-DELETE-FAIL'
	),
});

export {
	FolderList,
	FolderCreate,
	FolderDelete,

	type Folder,
	type FolderListResponse
}
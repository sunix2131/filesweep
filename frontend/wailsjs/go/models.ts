export namespace actions {
	
	export class Item {
	    id: string;
	    actionId: string;
	    sourcePath: string;
	    targetPath: string;
	    sourceSizeBytes: number;
	    sourceSha256: string;
	    status: string;
	    errorMessage: string;
	
	    static createFrom(source: any = {}) {
	        return new Item(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.actionId = source["actionId"];
	        this.sourcePath = source["sourcePath"];
	        this.targetPath = source["targetPath"];
	        this.sourceSizeBytes = source["sourceSizeBytes"];
	        this.sourceSha256 = source["sourceSha256"];
	        this.status = source["status"];
	        this.errorMessage = source["errorMessage"];
	    }
	}
	export class Action {
	    id: string;
	    actionType: string;
	    status: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    startedAt?: any;
	    // Go type: time
	    completedAt?: any;
	    undoAvailable: boolean;
	    // Go type: time
	    undoExpiresAt?: any;
	    summary: string;
	    errorMessage: string;
	    items?: Item[];
	
	    static createFrom(source: any = {}) {
	        return new Action(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.actionType = source["actionType"];
	        this.status = source["status"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.startedAt = this.convertValues(source["startedAt"], null);
	        this.completedAt = this.convertValues(source["completedAt"], null);
	        this.undoAvailable = source["undoAvailable"];
	        this.undoExpiresAt = this.convertValues(source["undoExpiresAt"], null);
	        this.summary = source["summary"];
	        this.errorMessage = source["errorMessage"];
	        this.items = this.convertValues(source["items"], Item);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace dto {
	
	export class BuildActionPlanRequest {
	    actionType: string;
	    items: actions.Item[];
	
	    static createFrom(source: any = {}) {
	        return new BuildActionPlanRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.actionType = source["actionType"];
	        this.items = this.convertValues(source["items"], actions.Item);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ExportResultDto {
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new ExportResultDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	    }
	}
	export class FolderSelectionResult {
	    paths: string[];
	    hasDangerousPath: boolean;
	    dangerousPathNote: string;
	
	    static createFrom(source: any = {}) {
	        return new FolderSelectionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.paths = source["paths"];
	        this.hasDangerousPath = source["hasDangerousPath"];
	        this.dangerousPathNote = source["dangerousPathNote"];
	    }
	}
	export class PageQuery {
	    scanId: string;
	    limit: number;
	    offset: number;
	    category: string;
	    search: string;
	    minSizeBytes: number;
	
	    static createFrom(source: any = {}) {
	        return new PageQuery(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.scanId = source["scanId"];
	        this.limit = source["limit"];
	        this.offset = source["offset"];
	        this.category = source["category"];
	        this.search = source["search"];
	        this.minSizeBytes = source["minSizeBytes"];
	    }
	}
	export class PaginatedDuplicateGroupsDto {
	    items: any;
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new PaginatedDuplicateGroupsDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = source["items"];
	        this.total = source["total"];
	    }
	}
	export class PaginatedFilesDto {
	    items: scan.File[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new PaginatedFilesDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], scan.File);
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StartScanRequest {
	    paths: string[];
	    confirmDangerousFolders: boolean;
	
	    static createFrom(source: any = {}) {
	        return new StartScanRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.paths = source["paths"];
	        this.confirmDangerousFolders = source["confirmDangerousFolders"];
	    }
	}

}

export namespace duplicates {
	
	export class Group {
	    id: string;
	    scanSessionId: string;
	    sha256: string;
	    fileSizeBytes: number;
	    filesCount: number;
	    estimatedReclaimableBytes: number;
	    recommendedFileId: string;
	    category: string;
	    members?: scan.File[];
	
	    static createFrom(source: any = {}) {
	        return new Group(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.scanSessionId = source["scanSessionId"];
	        this.sha256 = source["sha256"];
	        this.fileSizeBytes = source["fileSizeBytes"];
	        this.filesCount = source["filesCount"];
	        this.estimatedReclaimableBytes = source["estimatedReclaimableBytes"];
	        this.recommendedFileId = source["recommendedFileId"];
	        this.category = source["category"];
	        this.members = this.convertValues(source["members"], scan.File);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace events {
	
	export class ScanProgress {
	    scanId: string;
	    phase: string;
	    processedFiles: number;
	    totalFiles: number;
	    currentPath: string;
	    percent: number;
	
	    static createFrom(source: any = {}) {
	        return new ScanProgress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.scanId = source["scanId"];
	        this.phase = source["phase"];
	        this.processedFiles = source["processedFiles"];
	        this.totalFiles = source["totalFiles"];
	        this.currentPath = source["currentPath"];
	        this.percent = source["percent"];
	    }
	}

}

export namespace scan {
	
	export class File {
	    id: string;
	    scanSessionId: string;
	    rootId: string;
	    absolutePath: string;
	    relativePath: string;
	    name: string;
	    extension: string;
	    mimeType: string;
	    category: string;
	    sizeBytes: number;
	    // Go type: time
	    modifiedAt: any;
	    // Go type: time
	    createdAtIfAvailable?: any;
	    sha256: string;
	    isHidden: boolean;
	    isAccessible: boolean;
	    isSymlink: boolean;
	    isUnstable: boolean;
	    scanError: string;
	
	    static createFrom(source: any = {}) {
	        return new File(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.scanSessionId = source["scanSessionId"];
	        this.rootId = source["rootId"];
	        this.absolutePath = source["absolutePath"];
	        this.relativePath = source["relativePath"];
	        this.name = source["name"];
	        this.extension = source["extension"];
	        this.mimeType = source["mimeType"];
	        this.category = source["category"];
	        this.sizeBytes = source["sizeBytes"];
	        this.modifiedAt = this.convertValues(source["modifiedAt"], null);
	        this.createdAtIfAvailable = this.convertValues(source["createdAtIfAvailable"], null);
	        this.sha256 = source["sha256"];
	        this.isHidden = source["isHidden"];
	        this.isAccessible = source["isAccessible"];
	        this.isSymlink = source["isSymlink"];
	        this.isUnstable = source["isUnstable"];
	        this.scanError = source["scanError"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Session {
	    id: string;
	    status: string;
	    // Go type: time
	    startedAt: any;
	    // Go type: time
	    completedAt?: any;
	    // Go type: time
	    cancelledAt?: any;
	    selectedPaths: string[];
	    filesCount: number;
	    totalSizeBytes: number;
	    duplicateGroupsCount: number;
	    reclaimableBytes: number;
	    errorsCount: number;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Session(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.status = source["status"];
	        this.startedAt = this.convertValues(source["startedAt"], null);
	        this.completedAt = this.convertValues(source["completedAt"], null);
	        this.cancelledAt = this.convertValues(source["cancelledAt"], null);
	        this.selectedPaths = source["selectedPaths"];
	        this.filesCount = source["filesCount"];
	        this.totalSizeBytes = source["totalSizeBytes"];
	        this.duplicateGroupsCount = source["duplicateGroupsCount"];
	        this.reclaimableBytes = source["reclaimableBytes"];
	        this.errorsCount = source["errorsCount"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace settings {
	
	export class Settings {
	    language: string;
	    theme: string;
	    includeHiddenFiles: boolean;
	    followSymlinks: boolean;
	    maxHashWorkers: number;
	    maxPreviewFileSizeMb: number;
	    scanExcludedFolderNames: string[];
	    recentFolders: string[];
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.language = source["language"];
	        this.theme = source["theme"];
	        this.includeHiddenFiles = source["includeHiddenFiles"];
	        this.followSymlinks = source["followSymlinks"];
	        this.maxHashWorkers = source["maxHashWorkers"];
	        this.maxPreviewFileSizeMb = source["maxPreviewFileSizeMb"];
	        this.scanExcludedFolderNames = source["scanExcludedFolderNames"];
	        this.recentFolders = source["recentFolders"];
	    }
	}

}


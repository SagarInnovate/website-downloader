export namespace main {
	
	export class HistoryItem {
	    jobId: string;
	    url: string;
	    mode: string;
	    status: string;
	    zipPath: string;
	    // Go type: time
	    startTime: any;
	    // Go type: time
	    endTime: any;
	    pages: number;
	    assets: number;
	    totalSize: number;
	
	    static createFrom(source: any = {}) {
	        return new HistoryItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.jobId = source["jobId"];
	        this.url = source["url"];
	        this.mode = source["mode"];
	        this.status = source["status"];
	        this.zipPath = source["zipPath"];
	        this.startTime = this.convertValues(source["startTime"], null);
	        this.endTime = this.convertValues(source["endTime"], null);
	        this.pages = source["pages"];
	        this.assets = source["assets"];
	        this.totalSize = source["totalSize"];
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
	export class ScrapeResponse {
	    jobId: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new ScrapeResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.jobId = source["jobId"];
	        this.status = source["status"];
	    }
	}

}

export namespace models {
	
	export class JobStatus {
	    jobId: string;
	    status: string;
	    url: string;
	    pagesDiscovered: number;
	    pagesDownloaded: number;
	    assetsDownloaded: number;
	    totalSize: number;
	    currentPage: string;
	    progress: number;
	    error?: string;
	    // Go type: time
	    startTime: any;
	    // Go type: time
	    endTime?: any;
	
	    static createFrom(source: any = {}) {
	        return new JobStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.jobId = source["jobId"];
	        this.status = source["status"];
	        this.url = source["url"];
	        this.pagesDiscovered = source["pagesDiscovered"];
	        this.pagesDownloaded = source["pagesDownloaded"];
	        this.assetsDownloaded = source["assetsDownloaded"];
	        this.totalSize = source["totalSize"];
	        this.currentPage = source["currentPage"];
	        this.progress = source["progress"];
	        this.error = source["error"];
	        this.startTime = this.convertValues(source["startTime"], null);
	        this.endTime = this.convertValues(source["endTime"], null);
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


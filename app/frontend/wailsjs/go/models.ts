export namespace types {
	
	export class Connect {
	    id: number;
	    name: string;
	    bootstrap_servers: string;
	    tls: string;
	    skipTLSVerify: string;
	    tls_cert_file: string;
	    tls_key_file: string;
	    tls_ca_file: string;
	    sasl: string;
	    sasl_mechanism: string;
	    sasl_user: string;
	    sasl_pwd: string;
	
	    static createFrom(source: any = {}) {
	        return new Connect(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.bootstrap_servers = source["bootstrap_servers"];
	        this.tls = source["tls"];
	        this.skipTLSVerify = source["skipTLSVerify"];
	        this.tls_cert_file = source["tls_cert_file"];
	        this.tls_key_file = source["tls_key_file"];
	        this.tls_ca_file = source["tls_ca_file"];
	        this.sasl = source["sasl"];
	        this.sasl_mechanism = source["sasl_mechanism"];
	        this.sasl_user = source["sasl_user"];
	        this.sasl_pwd = source["sasl_pwd"];
	    }
	}
	export class Config {
	    width: number;
	    height: number;
	    language: string;
	    theme: string;
	    connects: Connect[];
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.width = source["width"];
	        this.height = source["height"];
	        this.language = source["language"];
	        this.theme = source["theme"];
	        this.connects = this.convertValues(source["connects"], Connect);
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
	
	export class ResultResp {
	    result: {[key: string]: any};
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new ResultResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.result = source["result"];
	        this.err = source["err"];
	    }
	}
	export class ResultsResp {
	    results: any[];
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new ResultsResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.results = source["results"];
	        this.err = source["err"];
	    }
	}
	export class Tag {
	    tag_name: string;
	    body: string;
	
	    static createFrom(source: any = {}) {
	        return new Tag(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tag_name = source["tag_name"];
	        this.body = source["body"];
	    }
	}

}


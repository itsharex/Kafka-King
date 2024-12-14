export namespace frontend {
	
	export class FileFilter {
	    DisplayName: string;
	    Pattern: string;
	
	    static createFrom(source: any = {}) {
	        return new FileFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DisplayName = source["DisplayName"];
	        this.Pattern = source["Pattern"];
	    }
	}
	export class OpenDialogOptions {
	    DefaultDirectory: string;
	    DefaultFilename: string;
	    Title: string;
	    Filters: FileFilter[];
	    ShowHiddenFiles: boolean;
	    CanCreateDirectories: boolean;
	    ResolvesAliases: boolean;
	    TreatPackagesAsDirectories: boolean;
	
	    static createFrom(source: any = {}) {
	        return new OpenDialogOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DefaultDirectory = source["DefaultDirectory"];
	        this.DefaultFilename = source["DefaultFilename"];
	        this.Title = source["Title"];
	        this.Filters = this.convertValues(source["Filters"], FileFilter);
	        this.ShowHiddenFiles = source["ShowHiddenFiles"];
	        this.CanCreateDirectories = source["CanCreateDirectories"];
	        this.ResolvesAliases = source["ResolvesAliases"];
	        this.TreatPackagesAsDirectories = source["TreatPackagesAsDirectories"];
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
	    kerberos_user_keytab: string;
	    kerberos_krb5_conf: string;
	    Kerberos_user: string;
	    Kerberos_realm: string;
	    kerberos_service_name: string;
	
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
	        this.kerberos_user_keytab = source["kerberos_user_keytab"];
	        this.kerberos_krb5_conf = source["kerberos_krb5_conf"];
	        this.Kerberos_user = source["Kerberos_user"];
	        this.Kerberos_realm = source["Kerberos_realm"];
	        this.kerberos_service_name = source["kerberos_service_name"];
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


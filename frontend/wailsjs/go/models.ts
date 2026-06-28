export namespace main {
	
	export class NewsItem {
	    date: string;
	    title: string;
	    content: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new NewsItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.title = source["title"];
	        this.content = source["content"];
	        this.url = source["url"];
	    }
	}
	export class ServerStatusInfo {
	    online: boolean;
	    playersOnline: number;
	    totalPlayers: number;
	
	    static createFrom(source: any = {}) {
	        return new ServerStatusInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.online = source["online"];
	        this.playersOnline = source["playersOnline"];
	        this.totalPlayers = source["totalPlayers"];
	    }
	}

}


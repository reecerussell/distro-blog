import { Component, OnInit } from "@angular/core";
import { ApiService } from "../api.service";

@Component({
    selector: "app-page-list",
    templateUrl: "./page-list.component.html",
    styleUrls: ["./page-list.component.scss"],
})
export class PageListComponent implements OnInit {
    pages: any = null;
    searchTerm: string = "";
    error: string = null;
    loading: boolean = false;

    constructor(private api: ApiService) {}

    async ngOnInit() {
        await this.loadPages();
    }

    async loadPages(): Promise<void> {
        if (this.loading) {
            return;
        }

        this.loading = true;

        const res = await this.api.Pages.List();
        if (res.ok) {
            this.pages = res.data;
        } else {
            this.error = res.error;
        }

        this.loading = false;
    }

    getPages(): any {
        const term = this.searchTerm.replace(/ +/g, " ").toLowerCase();
        if (term.length < 1 || !this.pages) {
            return this.pages;
        }

        return this.pages.filter((p) => {
            const title = p.title.replace(/\s+/g, " ").toLowerCase();
            const description = p.description
                .replace(/\s+/g, " ")
                .toLowerCase();

            return title.indexOf(term) > -1 || description.indexOf(term) > -1;
        });
    }
}

import { Component, OnInit } from "@angular/core";
import { ApiService } from "../api.service";
import { Page } from "../models/page";

@Component({
    selector: "app-blog-list",
    templateUrl: "./blog-list.component.html",
    styleUrls: ["./blog-list.component.scss"],
})
export class BlogListComponent implements OnInit {
    blogs: Page[] = null;
    searchTerm: string = "";
    error: string = null;
    loading: boolean = false;

    constructor(private api: ApiService) {}

    async ngOnInit() {
        await this.loadBlogs();
    }

    async loadBlogs(): Promise<void> {
        if (this.loading) {
            return;
        }

        this.loading = true;

        const res = await this.api.Blogs.List();
        if (res.ok) {
            this.blogs = res.data;
        } else {
            this.error = res.error;
        }

        this.loading = false;
    }

    getBlogs(): any {
        const term = this.searchTerm.replace(/ +/g, " ").toLowerCase();
        if (term.length < 1 || !this.blogs) {
            return this.blogs;
        }

        return this.blogs.filter((p) => {
            const title = p.title.replace(/\s+/g, " ").toLowerCase();
            const description = p.description
                .replace(/\s+/g, " ")
                .toLowerCase();

            return title.indexOf(term) > -1 || description.indexOf(term) > -1;
        });
    }
}

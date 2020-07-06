import { Component, OnInit } from "@angular/core";
import { ApiService } from "../api.service";
import { ActivatedRoute } from "@angular/router";
import { Title } from "@angular/platform-browser";
import { Page } from "../models/page";

@Component({
    selector: "app-page-info",
    templateUrl: "./page-info.component.html",
    styleUrls: ["./page-info.component.scss"],
})
export class PageInfoComponent implements OnInit {
    model: Page;
    error: string = null;
    loading: boolean = false;

    constructor(
        private api: ApiService,
        private titleService: Title,
        private route: ActivatedRoute
    ) {}

    async ngOnInit(): Promise<void> {
        this.titleService.setTitle("Page Info - Distro Blog Admin");

        this.route.paramMap.subscribe(
            async (params) => await this.fetchPage(params.get("id"))
        );
    }

    async fetchPage(id: string) {
        if (this.loading) {
            return;
        }

        this.loading = true;

        const res = await this.api.Pages.Get(id, "audit");
        if (res.ok) {
            this.error = null;
            this.model = res.data as Page;

            if (this.model.isBlog) {
                this.titleService.setTitle("Blog Info - Distro Blog Admin");
            }
        } else {
            this.error = res.error;
        }

        this.loading = false;
    }

    getAuditMessage(message): string {
        const pageType = this.model.isBlog ? "Blog" : "Page";

        switch (message) {
            case "PAGE_UPDATED":
                return pageType + " was updated.";
            case "PAGE_CREATED":
                return pageType + " was created.";
            case "PAGE_DEACTIVATED":
                return pageType + " was deactivated.";
            case "PAGE_ACTIVATED":
                return pageType + " was activated.";
            default:
                return message;
        }
    }

    getMessageColorClass(message): string {
        switch (message) {
            case "PAGE_UPDATED":
                return "message-updated";
            case "PAGE_CREATED":
                return "message-created";
            case "PAGE_DEACTIVATED":
                return "message-deactivated";
            case "PAGE_ACTIVATED":
                return "message-activated";
            default:
                return null;
        }
    }
}

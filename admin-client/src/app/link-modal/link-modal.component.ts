import { Component, OnInit, Input } from "@angular/core";
import { ApiService } from "../api.service";
import Link from "../models/link";

@Component({
    selector: "app-link-modal",
    templateUrl: "./link-modal.component.html",
    styleUrls: ["./link-modal.component.scss"],
})
export class LinkModalComponent implements OnInit {
    @Input("link")
    link: Link;

    @Input("open")
    isOpen: boolean = true;

    isPageDropdownOpen: boolean = false;
    isPageOptionsLoading: boolean = false;
    pageOptions: any[];
    selectedPageId: string = null;

    constructor(private api: ApiService) {}

    async ngOnInit(): Promise<void> {
        if (!this.link) {
            this.link = {
                text: "",
                target: "_self",
                pageId: null,
                url: "",
            };
        }
    }

    async fetchPageOptions(): Promise<void> {
        if (this.isPageOptionsLoading) {
            return;
        }

        this.isPageOptionsLoading = true;

        const { ok, error, data } = await this.api.Pages.DropdownOptions();
        if (!ok) {
            console.log(error);
        } else {
            this.pageOptions = data;
        }

        this.isPageOptionsLoading = false;
    }

    toggle() {
        this.isOpen = !this.isOpen;
    }

    openPageDropdown() {
        this.isPageDropdownOpen = true;
        if (!this.pageOptions) {
            this.fetchPageOptions();
        }
    }

    closePageDropdown() {
        this.isPageDropdownOpen = false;
    }

    selectPage(id: string, url: string) {
        if (this.selectedPageId === id) {
            this.selectedPageId = null;
        } else {
            this.selectedPageId = id;
            this.link.url = url;
        }
    }
}

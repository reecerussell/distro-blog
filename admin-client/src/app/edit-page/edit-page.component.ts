import { Component, OnInit, ViewChild, ElementRef } from "@angular/core";
import { Page } from "../models/page";
import { ApiService } from "../api.service";
import { ActivatedRoute } from "@angular/router";
import { Title } from "@angular/platform-browser";
import * as ClassicEditor from "@ckeditor/ckeditor5-build-classic";

@Component({
    selector: "app-edit-page",
    templateUrl: "./edit-page.component.html",
    styleUrls: ["./edit-page.component.scss"],
})
export class EditPageComponent implements OnInit {
    @ViewChild("imageUpload") imageUpload: ElementRef;

    model: Page;
    contentEditor = ClassicEditor;
    error: string = null;
    loading: boolean = false;

    titleError: string = null;
    descriptionError: string = null;

    constructor(
        private api: ApiService,
        private titleService: Title,
        private route: ActivatedRoute
    ) {}

    ngOnInit(): void {
        this.titleService.setTitle("Edit Page - Distro Blog Admin");

        this.route.paramMap.subscribe(
            async (params) => await this.fetchPage(params.get("id"))
        );
    }

    async fetchPage(id: string) {
        const res = await this.api.Pages.Get(id);
        if (res.ok) {
            this.error = null;
            this.model = res.data as Page;

            if (this.model.isBlog) {
                this.titleService.setTitle("Edit Blog - Distro Blog Admin");
            }
        } else {
            this.error = res.error;
        }
    }

    validateTitle(): any {
        const title = this.model.title;
        if (title.length < 1) {
            this.titleError = "Please enter a title..";
            return false;
        } else if (title.length > 255) {
            this.titleError = "Title cannot be longer than 255 characters";
            return false;
        } else {
            this.titleError = null;
        }

        return true;
    }

    validateDescription(): any {
        const description = this.model.description;
        if (description.length < 1) {
            this.descriptionError = "Please enter a description..";
            return false;
        } else if (description.length > 255) {
            this.descriptionError =
                "Description cannot be longer than 255 characters";
            return false;
        } else {
            this.descriptionError = null;
        }

        return true;
    }

    validate() {
        let valid = true;

        if (!this.validateTitle()) {
            valid = false;
        }

        if (!this.validateDescription()) {
            valid = false;
        }

        return valid;
    }

    async onSubmit() {
        if (this.loading || !this.validate()) {
            return;
        }

        this.loading = true;

        let imageBlog: Blob = null;
        const fileUpload = this.imageUpload.nativeElement;
        if (fileUpload.files && fileUpload.files.length > 0) {
            const buffer = await fileUpload.files[0].arrayBuffer();
            imageBlog = new Blob([
                new Uint8Array(buffer, 0, buffer.byteLength),
            ]);
        }

        const res = await this.api.Pages.Update(this.model, imageBlog);
        this.error = res.ok ? null : res.error;

        this.loading = false;
    }
}

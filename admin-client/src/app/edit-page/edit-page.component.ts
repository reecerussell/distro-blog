import { Component, OnInit, ViewChild, ElementRef } from "@angular/core";
import { Page } from "../models/page";
import { ApiService } from "../api.service";
import { ActivatedRoute } from "@angular/router";
import { Title } from "@angular/platform-browser";
import * as ClassicEditor from "@ckeditor/ckeditor5-build-classic";

const allowedUrlChars = "abcdefghijklmnopqrstuvwxyz1234567890-";

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
    fileInputLabel: string = "Choose file";

    titleError: string = null;
    descriptionError: string = null;
    urlError: string = null;
    seoTitleError: string = null;
    seoDescriptionError: string = null;

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
            this.descriptionError = "Please enter a description.";
            return false;
        } else if (description.length > 255) {
            this.descriptionError =
                "Description cannot be longer than 255 characters.";
            return false;
        } else {
            this.descriptionError = null;
        }

        return true;
    }

    validateUrl(): any {
        const url = this.model.url;
        if (url.length > 255) {
            this.urlError = "URL cannot be longer than 255 characters.";
            return false;
        } else {
            const chars = allowedUrlChars.split("");
            const urlChars = url.toLowerCase().split("");
            for (let i = 0; i < urlChars.length; i++) {
                const char = urlChars[i];
                let valid = false;

                for (let j = 0; j < chars.length; j++) {
                    if (char === chars[j]) {
                        valid = true;
                        break;
                    }
                }

                if (!valid) {
                    this.urlError = `The character '${char}' is not allowed in URLs.`;
                    return false;
                }
            }

            this.urlError = null;
        }

        return true;
    }

    validateSeoTitle(): any {
        const title = this.model.seo?.title;
        if (title.length < 1) {
            this.seoTitleError = "Please enter a title.";
        } else if (title.length > 255) {
            this.seoTitleError =
                "Title cannot be greater than 255 characters long.";
        } else {
            this.seoTitleError = null;
        }

        return this.seoTitleError == null;
    }

    validateSeoDescription(): any {
        const description = this.model.seo?.description;
        if (description.length < 1) {
            this.seoDescriptionError = "Please enter a description.";
        } else if (description.length > 255) {
            this.seoDescriptionError =
                "Description cannot be greater than 255 characters long.";
        } else {
            this.seoDescriptionError = null;
        }

        return this.seoDescriptionError == null;
    }

    validate() {
        let valid = true;

        if (!this.validateTitle()) {
            valid = false;
        }

        if (!this.validateDescription()) {
            valid = false;
        }

        if (!this.validateUrl()) {
            valid = false;
        }

        if (!this.validateSeoTitle()) {
            valid = false;
        }

        if (!this.validateSeoDescription()) {
            valid = false;
        }

        return valid;
    }

    onFileChange() {
        const upload = this.imageUpload.nativeElement;
        if (upload.files.length > 0) {
            this.fileInputLabel = upload.files[0].name;
        } else {
            this.fileInputLabel = "Choose image";
        }
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
        if (!res.ok) {
            this.error = res.error;
        } else {
            this.error = null;
            await this.fetchPage(this.model.id);
        }

        this.loading = false;
    }
}

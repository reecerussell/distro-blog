import { Component, OnInit } from "@angular/core";
import { CreatePage } from "../models/page";
import { ApiService } from "../api.service";
import { Router } from "@angular/router";
import { Title } from "@angular/platform-browser";
import * as ClassicEditor from "@ckeditor/ckeditor5-build-classic";

const allowedUrlChars = "abcdefghijklmnopqrstuvwxyz1234567890-";

@Component({
    selector: "app-create-page",
    templateUrl: "./create-page.component.html",
    styleUrls: ["./create-page.component.scss"],
})
export class CreatePageComponent implements OnInit {
    model: CreatePage;
    error: string = null;
    loading: boolean = false;
    contentEditor = ClassicEditor;
    isBlog: boolean = false;
    isUrlDirty: boolean = false;

    titleError: string = null;
    descriptionError: string = null;
    urlError: string = null;

    constructor(
        private api: ApiService,
        private router: Router,
        private titleService: Title
    ) {}

    ngOnInit(): void {
        this.titleService.setTitle("New Page - Distro Blog Admin");

        this.model = {
            title: "",
            description: "",
            content: "",
            url: "",
        };
    }

    onTitleChange(): any {
        this.validateTitle();

        if (!this.isUrlDirty) {
            setTimeout(() => {
                const chars = allowedUrlChars.split("");
                let url = this.model.title.toLowerCase().replace(" ", "-");
                let urlChars = url.split("");

                let count = urlChars.length;
                let i = 0;

                while (i < count) {
                    if (urlChars[i] === " ") {
                        urlChars[i] = "-";
                    } else if (chars.indexOf(urlChars[i]) == -1) {
                        urlChars.splice(i, 1);
                        count--;
                    }

                    i++;
                }

                this.model.url = urlChars.join("");
            }, 200);
        }
    }

    onUrlChange(): any {
        this.validateUrl();
        this.isUrlDirty = true;
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

        return valid;
    }

    async onSubmit() {
        if (this.loading || !this.validate()) {
            return;
        }

        this.loading = true;

        const res = await this.api.Pages.Create(this.model);
        if (res.ok) {
            this.router.navigateByUrl("pages/" + res.data);
        } else {
            this.error = res.error;
        }

        this.loading = false;
    }
}

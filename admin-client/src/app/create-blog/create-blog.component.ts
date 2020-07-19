import { Component, OnInit } from "@angular/core";
import { CreatePage } from "../models/page";
import { ApiService } from "../api.service";
import { Router } from "@angular/router";
import { Title } from "@angular/platform-browser";
import * as ClassicEditor from "@ckeditor/ckeditor5-build-classic";

const allowedUrlChars = "abcdefghijklmnopqrstuvwxyz1234567890-";

@Component({
    selector: "app-create-blog",
    templateUrl: "./create-blog.component.html",
    styleUrls: ["./create-blog.component.scss"],
})
export class CreateBlogComponent implements OnInit {
    model: CreatePage;
    error: string = null;
    loading: boolean = false;
    contentEditor = ClassicEditor;
    isBlog: boolean = false;
    isUrlDirty: boolean = false;
    isSeoTitleDirty: boolean = false;
    isSeoDescriptionDirty: boolean = false;

    titleError: string = null;
    descriptionError: string = null;
    urlError: string = null;
    seoTitleError: string = null;
    seoDescriptionError: string = null;

    constructor(
        private api: ApiService,
        private router: Router,
        private titleService: Title
    ) {}

    ngOnInit(): void {
        this.titleService.setTitle("New Blog - Distro Blog Admin");

        this.model = {
            title: "",
            description: "",
            content: "",
            url: "",
            seo: {
                title: "",
                description: "",
                index: false,
                follow: false,
            },
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

        if (!this.isSeoTitleDirty) {
            this.model.seo.title = this.model.title;
        }
    }

    onUrlChange(): any {
        this.validateUrl();
        this.isUrlDirty = true;
    }

    onDescriptionChange(): any {
        this.validateDescription();

        if (!this.isSeoDescriptionDirty) {
            this.model.seo.description = this.model.description;
        }
    }

    onSeoTitleChange(): any {
        this.validateSeoTitle();

        this.isSeoTitleDirty = true;
    }

    onSeoDescriptionChange(): any {
        this.validateSeoDescription();

        this.isSeoDescriptionDirty = true;
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

    validateSeoTitle(): any {
        const title = this.model.seo?.title;
        if (title.length > 255) {
            this.seoTitleError =
                "Title cannot be greater than 255 characters long.";
        } else {
            this.seoTitleError = null;
        }

        return this.seoTitleError == null;
    }

    validateSeoDescription(): any {
        const description = this.model.seo?.description;
        if (description.length > 255) {
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
            return false;
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

    async onSubmit() {
        if (this.loading || !this.validate()) {
            return;
        }

        this.loading = true;

        const res = await this.api.Blogs.Create(this.model);
        if (res.ok) {
            this.router.navigateByUrl("blogs/" + res.data);
        } else {
            this.error = res.error;
        }

        this.loading = false;
    }
}

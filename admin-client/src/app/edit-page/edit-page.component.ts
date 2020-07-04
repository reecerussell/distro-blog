import { Component, OnInit, AfterViewInit, ViewChild } from "@angular/core";
import { UpdatePage } from "../models/page";
import { ApiService } from "../api.service";
import { Router, ActivatedRoute } from "@angular/router";
import { Title } from "@angular/platform-browser";
import * as ClassicEditor from "@ckeditor/ckeditor5-build-classic";
import { NgForm } from "@angular/forms";

@Component({
    selector: "app-edit-page",
    templateUrl: "./edit-page.component.html",
    styleUrls: ["./edit-page.component.scss"],
})
export class EditPageComponent implements OnInit {
    public model: UpdatePage;
    public contentEditor = ClassicEditor;

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
            this.model = res.data as UpdatePage;
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

        const res = await this.api.Pages.Update(this.model);
        this.error = res.ok ? null : res.error;

        this.loading = false;
    }
}

import { Component, OnInit } from "@angular/core";
import Setting from "../models/setting";
import { ApiService } from "../api.service";
import { ActivatedRoute } from "@angular/router";
import { Title } from "@angular/platform-browser";

@Component({
    selector: "app-edit-setting",
    templateUrl: "./edit-setting.component.html",
    styleUrls: ["./edit-setting.component.scss"],
})
export class EditSettingComponent implements OnInit {
    model: Setting;
    error: string = null;
    loading: boolean = false;
    valueError: string = null;

    constructor(
        private api: ApiService,
        private titleService: Title,
        private route: ActivatedRoute
    ) {}

    async ngOnInit() {
        this.titleService.setTitle("Edit Setting - Distro Blog Admin");

        this.route.paramMap.subscribe(
            async (params) => await this.fetchSetting(params.get("key"))
        );
    }

    async fetchSetting(key: string) {
        const res = await this.api.Settings.Get(key);
        if (res.ok) {
            this.error = null;
            this.model = res.data as Setting;
        } else {
            this.error = res.error;
        }
    }

    validateValue(): any {
        if (this.model.value.length > 255) {
            this.valueError =
                "Value cannot be greater than 255 characters long.";
            return false;
        } else {
            this.validate;
        }

        return true;
    }

    validate() {
        let valid = true;

        if (!this.validateValue()) {
            valid = false;
        }

        return valid;
    }

    async onSubmit() {
        if (this.loading || !this.validate()) {
            return;
        }

        this.loading = true;

        const res = await this.api.Settings.Update(this.model);
        this.error = res.ok ? null : res.error;

        this.loading = false;
    }

    getHelpText(): string {
        switch (this.model?.key) {
            case "TITLE_FORMAT":
                return 'The format of page titles. Use replacements "{TITLE}" for a page\'s title and "{SITE_NAME}" for the site name.';
            default:
                return null;
        }
    }
}

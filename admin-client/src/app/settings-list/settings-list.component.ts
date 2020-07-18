import { Component, OnInit } from "@angular/core";
import { ApiService } from "../api.service";
import Setting from "../models/setting";

@Component({
    selector: "app-settings-list",
    templateUrl: "./settings-list.component.html",
    styleUrls: ["./settings-list.component.scss"],
})
export class SettingsListComponent implements OnInit {
    settings: Setting[] = null;
    searchTerm: string = "";
    error: string = null;
    loading: boolean = false;

    constructor(private api: ApiService) {}

    async ngOnInit() {
        await this.loadSettings();
    }

    async loadSettings(): Promise<void> {
        if (this.loading) {
            return;
        }

        this.loading = true;

        const res = await this.api.Settings.List();
        if (res.ok) {
            this.settings = res.data;
        } else {
            this.error = res.error;
        }

        this.loading = false;
    }

    getSettings(): any {
        const term = this.searchTerm.replace(/ +/g, " ").toLowerCase();
        if (term.length < 1 || !this.settings) {
            return this.settings;
        }

        return this.settings.filter((u) => {
            const key = u.key.replace(/\s+/g, " ").toLowerCase();
            const value = u.value.replace(/\s+/g, " ").toLowerCase();

            return key.indexOf(term) > -1 || value.indexOf(term) > -1;
        });
    }
}

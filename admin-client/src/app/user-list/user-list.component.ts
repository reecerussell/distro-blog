import { Component, OnInit } from "@angular/core";
import { ApiService } from "../api.service";

@Component({
    selector: "app-user-list",
    templateUrl: "./user-list.component.html",
    styleUrls: ["./user-list.component.scss"],
})
export class UserListComponent implements OnInit {
    users: any = null;
    searchTerm: string = "";
    error: string = null;
    loading: boolean = false;

    constructor(private api: ApiService) {}

    async ngOnInit() {
        await this.loadUsers();
    }

    async loadUsers(): Promise<void> {
        if (this.loading) {
            return;
        }

        this.loading = true;

        const res = await this.api.Users.List();
        if (res.ok) {
            this.users = res.data;
        } else {
            this.error = res.error;
        }

        this.loading = false;
    }

    getUsers(): any {
        const term = this.searchTerm.replace(/ +/g, " ").toLowerCase();
        if (term.length < 1 || !this.users) {
            return this.users;
        }

        return this.users.filter((u) => {
            const name = u.name.replace(/\s+/g, " ").toLowerCase();
            const email = u.email.replace(/\s+/g, " ").toLowerCase();

            return name.indexOf(term) > -1 || email.indexOf(term) > -1;
        });
    }
}

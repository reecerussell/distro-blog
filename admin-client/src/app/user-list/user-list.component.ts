import { Component, OnInit } from "@angular/core";
import { ApiService } from "../api.service";

@Component({
    selector: "app-user-list",
    templateUrl: "./user-list.component.html",
    styleUrls: ["./user-list.component.scss"],
})
export class UserListComponent implements OnInit {
    users: any = null;

    constructor(private api: ApiService) {}

    async ngOnInit() {
        const res = await this.api.Users.List();
        if (res.ok) {
            this.users = res.data;
        }
    }
}

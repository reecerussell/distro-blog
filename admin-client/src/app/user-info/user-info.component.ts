import { Component, OnInit } from "@angular/core";
import { ApiService } from "../api.service";
import { Router, ActivatedRoute } from "@angular/router";
import { Title } from "@angular/platform-browser";
import { User } from "../models/user";

@Component({
    selector: "app-user-info",
    templateUrl: "./user-info.component.html",
    styleUrls: ["./user-info.component.scss"],
})
export class UserInfoComponent implements OnInit {
    model: User;
    error: string = null;
    loading: boolean = false;
    openAudits: Map<number, boolean> = new Map<number, boolean>();

    constructor(
        private api: ApiService,
        private router: Router,
        private titleService: Title,
        private route: ActivatedRoute
    ) {}

    async ngOnInit(): Promise<void> {
        this.titleService.setTitle("User Info - Distro Blog Admin");

        this.route.paramMap.subscribe(
            async (params) => await this.fetchUser(params.get("id"))
        );
    }

    async fetchUser(id: string) {
        if (this.loading) {
            return;
        }

        this.loading = true;

        const res = await this.api.Users.Get(id, "audit");
        if (res.ok) {
            this.error = null;
            this.model = res.data as User;

            if (this.model.audit) {
                this.model.audit = this.model.audit.map((x) => ({
                    message: x.message,
                    userId: x.userId,
                    userFullname: x.userFullname,
                    date: new Date(x.date).toLocaleString(),
                    state: x.state,
                }));
            }
        } else {
            this.error = res.error;
        }

        this.loading = false;
    }

    toggleAuditDetails(id: number) {
        this.openAudits.set(id, !this.openAudits.get(id));
    }
}

import { Component, OnInit, Input } from "@angular/core";
import { ApiService } from "../api.service";

@Component({
    selector: "app-reset-password",
    templateUrl: "./reset-password.component.html",
    styleUrls: ["./reset-password.component.scss"],
})
export class ResetPasswordComponent implements OnInit {
    @Input()
    userId: string = "";

    loading: boolean = false;
    error: string = null;

    newPassword: string = null;

    constructor(private api: ApiService) {}

    ngOnInit(): void {}

    async onSubmit() {
        if (this.loading) {
            return;
        }

        this.loading = true;

        const res = await this.api.Users.ResetPassword(this.userId);
        if (res.ok) {
            this.error = null;
            this.newPassword = res.data;
        } else {
            this.newPassword = null;
            this.error = res.error;
        }

        this.loading = false;
    }
}

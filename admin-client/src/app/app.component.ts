import { Component } from "@angular/core";
import { UserService } from "./user.service";

@Component({
    selector: "app-root",
    templateUrl: "./app.component.html",
    styleUrls: ["./app.component.scss"],
})
export class AppComponent {
    title = "admin-client";

    isLoggedIn: boolean = false;

    constructor(private user: UserService) {
        this.isLoggedIn = user.IsAuthenticated();
    }

    ngOnInit() {
        setInterval(() => {
            if (!this.user.IsAuthenticated() && this.isLoggedIn) {
                this.user.Logout();
            }
        }, 5000);

        this.user.Listen(
            "app",
            () => (this.isLoggedIn = this.user.IsAuthenticated())
        );
    }
}

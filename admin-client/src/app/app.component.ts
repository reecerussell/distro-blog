import { Component, OnInit, OnDestroy } from "@angular/core";
import { UserService } from "./user.service";
import { Router } from "@angular/router";

@Component({
    selector: "app-root",
    templateUrl: "./app.component.html",
    styleUrls: ["./app.component.scss"],
})
export class AppComponent implements OnInit, OnDestroy {
    title = "admin-client";

    isLoggedIn: boolean = false;
    showMenu: boolean = false;

    constructor(private user: UserService, private router: Router) {
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

        this.router.events.subscribe(() => (this.showMenu = false));
    }

    ngOnDestroy() {
        this.user.Unlisten("app");
    }

    logout() {
        this.user.Logout();
    }

    toggleMenu() {
        this.showMenu = !this.showMenu;
    }
}

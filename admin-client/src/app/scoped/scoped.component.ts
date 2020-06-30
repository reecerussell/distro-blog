import { Component, OnInit } from "@angular/core";
import { UserService } from "../user.service";

@Component({
    selector: "app-scoped",
    templateUrl: "./scoped.component.html",
    styleUrls: ["./scoped.component.scss"],
})
export class ScopedComponent implements OnInit {
    display: boolean = false;

    constructor(public scope: string, private user: UserService) {
        this.display = user.IsInScope(scope);
    }

    ngOnInit(): void {}
}

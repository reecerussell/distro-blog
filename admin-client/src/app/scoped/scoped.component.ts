import { Component, OnInit, Input } from "@angular/core";
import { UserService } from "../user.service";

@Component({
    selector: "app-scoped",
    templateUrl: "./scoped.component.html",
    styleUrls: ["./scoped.component.scss"],
})
export class ScopedComponent implements OnInit {
    display: boolean = false;

    @Input()
    scope: string;

    constructor(private user: UserService) {}

    ngOnInit(): void {
        this.display = this.user.IsInScope(this.scope);
    }
}

import { Component, OnInit } from "@angular/core";
import { Router, ActivatedRoute } from "@angular/router";
import { ApiService } from "../api.service";
import { UpdateUser } from "../models/user";
import { Title } from "@angular/platform-browser";
import { switchMap } from "rxjs/operators";

@Component({
    selector: "app-edit-user",
    templateUrl: "./edit-user.component.html",
    styleUrls: ["./edit-user.component.scss"],
})
export class EditUserComponent implements OnInit {
    model: UpdateUser;
    error: string = null;
    loading: boolean = false;

    firstnameError: string = null;
    lastnameError: string = null;
    emailError: string = null;

    constructor(
        private api: ApiService,
        private router: Router,
        private titleService: Title,
        private route: ActivatedRoute
    ) {}

    async ngOnInit() {
        this.titleService.setTitle("Edit User - Distro Blog Admin");

        this.route.paramMap.subscribe(
            async (params) => await this.fetchUser(params.get("id"))
        );
    }

    async fetchUser(id: string) {
        const res = await this.api.Users.Get(id);
        if (res.ok) {
            this.error = null;
            this.model = res.data as UpdateUser;
        } else {
            this.error = res.error;
        }
    }

    validateFirstname(): any {
        const firstname = this.model.firstname;
        if (firstname.length < 1) {
            this.firstnameError = "Please enter the user's firtname.";
            return false;
        } else if (firstname.length > 45) {
            this.firstnameError =
                "Firstname cannot be longer than 45 characters";
            return false;
        } else {
            this.firstnameError = null;
        }

        return true;
    }

    validateLastname(): any {
        const lastname = this.model.lastname;
        if (lastname.length < 1) {
            this.lastnameError = "Please enter the user's lastname.";
            return false;
        } else if (lastname.length > 45) {
            this.lastnameError = "Lastname cannot be longer than 45 characters";
            return false;
        } else {
            this.lastnameError = null;
        }

        return true;
    }

    validateEmail(): any {
        const email = this.model.email;
        if (email.length < 1) {
            this.emailError = "This field is required.";
            return false;
        } else if (email.length > 100) {
            this.emailError = "Email cannot be longer than 100 characters";
            return false;
        } else if (
            !email.match("[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,6}")
        ) {
            this.emailError = "Please enter a valid email address.";
            return false;
        } else {
            this.emailError = null;
        }

        return true;
    }

    validate() {
        let valid = true;

        if (!this.validateFirstname()) {
            valid = false;
        }

        if (!this.validateLastname()) {
            valid = false;
        }

        if (!this.validateEmail()) {
            valid = false;
        }

        return valid;
    }

    async onSubmit() {
        if (this.loading || !this.validate()) {
            return;
        }

        this.loading = true;

        const res = await this.api.Users.Update(this.model);
        this.error = res.ok ? null : res.error;

        this.loading = false;
    }
}

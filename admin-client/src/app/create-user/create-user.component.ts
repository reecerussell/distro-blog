import { Component, OnInit } from "@angular/core";
import { Router } from "@angular/router";
import { ApiService } from "../api.service";
import { CreateUser } from "../models/user";
import { Title } from "@angular/platform-browser";

@Component({
    selector: "app-create-user",
    templateUrl: "./create-user.component.html",
    styleUrls: ["./create-user.component.scss"],
})
export class CreateUserComponent implements OnInit {
    model: CreateUser;
    error: string = null;
    loading: boolean = false;

    firstnameError: string = null;
    lastnameError: string = null;
    emailError: string = null;
    passwordError: string = null;

    constructor(
        private api: ApiService,
        private router: Router,
        private titleService: Title
    ) {}

    ngOnInit(): void {
        this.titleService.setTitle("New User - Distro Blog Admin");

        this.model = {
            firstname: "",
            lastname: "",
            email: "",
            password: "",
        };
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

    validatePassword(): any {
        const password = this.model.password;
        if (password.length < 1) {
            this.passwordError = "Please enter a password.";
            return false;
        } else {
            this.passwordError = null;
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

        if (!this.validatePassword()) {
            valid = false;
        }

        return valid;
    }

    async onSubmit() {
        if (this.loading || !this.validate()) {
            return;
        }

        this.loading = true;

        const res = await this.api.Users.Create(this.model);
        if (res.ok) {
            this.router.navigateByUrl("users");
        } else {
            this.error = res.error;
        }

        this.loading = false;
    }
}

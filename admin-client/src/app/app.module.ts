import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";
import { HttpClientModule } from "@angular/common/http";
import { FormsModule } from "@angular/forms";

import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { DashboardComponent } from "./dashboard/dashboard.component";
import { CreateUserComponent } from "./create-user/create-user.component";
import { UserListComponent } from "./user-list/user-list.component";
import { EditUserComponent } from "./edit-user/edit-user.component";
import { LoginComponent } from "./login/login.component";
import { UserInfoComponent } from "./user-info/user-info.component";
import { ScopedComponent } from './scoped/scoped.component';

@NgModule({
    declarations: [
        AppComponent,
        DashboardComponent,
        CreateUserComponent,
        UserListComponent,
        EditUserComponent,
        LoginComponent,
        UserInfoComponent,
        ScopedComponent,
    ],
    imports: [BrowserModule, AppRoutingModule, HttpClientModule, FormsModule],
    providers: [],
    bootstrap: [AppComponent],
})
export class AppModule {}

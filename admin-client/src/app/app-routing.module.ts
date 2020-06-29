import { NgModule } from "@angular/core";
import { Routes, RouterModule } from "@angular/router";
import { DashboardComponent } from "./dashboard/dashboard.component";
import { CreateUserComponent } from "./create-user/create-user.component";
import { UserListComponent } from "./user-list/user-list.component";
import { EditUserComponent } from "./edit-user/edit-user.component";
import { UserInfoComponent } from "./user-info/user-info.component";

const routes: Routes = [
    {
        path: "",
        component: DashboardComponent,
    },
    {
        path: "users",
        component: UserListComponent,
    },
    {
        path: "users/create",
        component: CreateUserComponent,
    },
    {
        path: "users/:id",
        component: EditUserComponent,
    },
    {
        path: "users/:id/info",
        component: UserInfoComponent,
    },
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule],
})
export class AppRoutingModule {}

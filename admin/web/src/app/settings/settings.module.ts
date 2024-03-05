import { ModuleWithProviders, NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { authGuard } from '../authentication/auth.guard';
import { SettingsComponent } from './settings/settings.component';
import { UiModule } from '../ui/ui.module';
import { LucideAngularModule } from 'lucide-angular';
import { ReactiveFormsModule } from '@angular/forms';
import { StatusComponent } from './status/status.component';

@NgModule({
  declarations: [SettingsComponent, StatusComponent],
  imports: [
    CommonModule,
    RouterModule.forChild([
      {
        path: 'config',
        redirectTo: 'config/settings'
      },
      {
        path: 'config/settings',
        component: SettingsComponent,
        canActivate: [authGuard]
      },
      {
        path: 'config/status',
        component: StatusComponent,
        canActivate: [authGuard]
      }
    ]),
    UiModule,
    LucideAngularModule,
    ReactiveFormsModule
  ]
})
export class SettingsModule {
  static forRoot(): ModuleWithProviders<SettingsModule> {
    return {
      ngModule: SettingsModule
    };
  }
}

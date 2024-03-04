import { ModuleWithProviders, NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { authGuard } from '../authentication/auth.guard';
import { SettingsComponent } from './settings/settings.component';
import { UiModule } from '../ui/ui.module';
import { LucideAngularModule } from 'lucide-angular';
import { ReactiveFormsModule } from '@angular/forms';

@NgModule({
  declarations: [SettingsComponent],
  imports: [
    CommonModule,
    RouterModule.forChild([
      {
        path: 'settings',
        component: SettingsComponent,
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

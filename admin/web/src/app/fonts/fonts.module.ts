import { ModuleWithProviders, NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FontListComponent } from './font-list/font-list.component';
import { authGuard } from '../authentication/auth.guard';
import { Check, LucideAngularModule, X } from 'lucide-angular';
import { FontDetailsComponent } from './font-details/font-details.component';

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    LucideAngularModule.pick({ Check, X }),
    RouterModule.forChild([
      {
        path: 'font',
        component: FontListComponent,
        canActivate: [authGuard]
      },
      {
        path: 'font/:fontName',
        component: FontDetailsComponent,
        canActivate: [authGuard]
      }
    ])
  ]
})
export class FontsModule {
  static forRoot(): ModuleWithProviders<FontsModule> {
    return {
      ngModule: FontsModule
    };
  }
}

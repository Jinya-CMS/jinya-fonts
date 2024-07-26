import { ModuleWithProviders, NgModule } from '@angular/core';
import { CommonModule, NgForOf, NgIf, NgSwitch, NgSwitchCase, NgSwitchDefault } from '@angular/common';
import { RouterLink, RouterModule } from '@angular/router';
import { FontListComponent } from './font-list/font-list.component';
import { Check, LucideAngularModule, X } from 'lucide-angular';
import { FontDetailsComponent } from './font-details/font-details.component';
import { UiModule } from '../ui/ui.module';
import { CreateFontDialogComponent } from './create-font-dialog/create-font-dialog.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { EditorModule } from '@tinymce/tinymce-angular';
import { authGuard } from '../authentication/auth.guard';
import { ConfirmComponent } from '../ui/confirm/confirm.component';
import { EditFontDialogComponent } from './edit-font-dialog/edit-font-dialog.component';
import { AddFileDialogComponent } from './add-file-dialog/add-file-dialog.component';
import { UpdateFileDialogComponent } from './update-file-dialog/update-file-dialog.component';

@NgModule({
  declarations: [
    FontListComponent,
    FontDetailsComponent,
    CreateFontDialogComponent,
    EditFontDialogComponent,
    AddFileDialogComponent,
    UpdateFileDialogComponent
  ],
  imports: [
    CommonModule,
    LucideAngularModule.pick({ Check, X }),
    NgForOf,
    NgIf,
    RouterLink,
    NgSwitchCase,
    NgSwitchDefault,
    NgSwitch,
    UiModule,
    ReactiveFormsModule,
    EditorModule,
    RouterModule.forChild([
      {
        path: 'font',
        redirectTo: 'font/all'
      },
      {
        path: 'font/:type',
        component: FontListComponent,
        canActivate: [authGuard]
      },
      {
        path: 'font/detail/:fontName',
        redirectTo: 'font/detail/:fontName/details'
      },
      {
        path: 'font/detail/:fontName/:section',
        component: FontDetailsComponent,
        canActivate: [authGuard]
      }
    ]),
    ConfirmComponent,
    FormsModule
  ]
})
export class FontsModule {
  static forRoot(): ModuleWithProviders<FontsModule> {
    return {
      ngModule: FontsModule
    };
  }
}

import { NgModule } from '@angular/core';
import { AppComponent } from './app.component';
import { NgIf } from '@angular/common';
import { RouterLink, RouterLinkActive, RouterModule } from '@angular/router';
import { BrowserModule } from '@angular/platform-browser';
import { OAuthModule } from 'angular-oauth2-oidc';
import { AuthenticationModule } from './authentication/authentication.module';
import { FontsModule } from './fonts/fonts.module';
import { ApiModule } from './api/api.module';
import { EditorModule, TINYMCE_SCRIPT_SRC } from '@tinymce/tinymce-angular';
import { SettingsModule } from './settings/settings.module';

@NgModule({
  declarations: [AppComponent],
  imports: [
    ApiModule.forRoot({ rootUrl: '' }),
    EditorModule,
    BrowserModule,
    NgIf,
    RouterLink,
    OAuthModule.forRoot(),
    AuthenticationModule.forRoot(),
    FontsModule.forRoot(),
    SettingsModule.forRoot(),
    RouterLinkActive,
    RouterModule.forRoot(
      [
        {
          path: '**',
          redirectTo: 'font/all'
        }
      ],
      {
        bindToComponentInputs: true
      }
    )
  ],
  bootstrap: [AppComponent],
  providers: [{ provide: TINYMCE_SCRIPT_SRC, useValue: 'tinymce/tinymce.min.js' }]
})
export class AppModule {}

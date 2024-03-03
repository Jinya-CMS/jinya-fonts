import { NgModule } from '@angular/core';
import { AppComponent } from './app.component';
import { NgIf } from '@angular/common';
import { RouterLink, RouterLinkActive, RouterModule } from '@angular/router';
import { BrowserModule } from '@angular/platform-browser';
import { OAuthModule } from 'angular-oauth2-oidc';
import { AuthenticationModule } from './authentication/authentication.module';
import { FontsModule } from './fonts/fonts.module';
import { ApiModule } from './api/api.module';

@NgModule({
  declarations: [AppComponent],
  imports: [
    ApiModule.forRoot({ rootUrl: '' }),
    BrowserModule,
    NgIf,
    RouterLink,
    OAuthModule.forRoot(),
    AuthenticationModule.forRoot(),
    FontsModule.forRoot(),
    RouterLinkActive,
    RouterModule.forRoot(
      [
        {
          path: '**',
          redirectTo: 'fonts'
        }
      ],
      {
        bindToComponentInputs: true
      }
    )
  ],
  bootstrap: [AppComponent]
})
export class AppModule {}

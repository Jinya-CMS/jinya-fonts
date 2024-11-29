import { ModuleWithProviders, NgModule, Optional, SkipSelf, inject, provideAppInitializer } from '@angular/core';
import { CommonModule } from '@angular/common';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import { AuthConfig, OAuthModule, OAuthModuleConfig, OAuthStorage } from 'angular-oauth2-oidc';
import { AuthenticationService } from './authentication.service';
import { authAppInitializerFactory } from './initializer.factory';
import { authConfig } from './authConfig';
import { authModuleConfig } from './authModuleConfig';
import { RouterModule } from '@angular/router';
import { LoginComponent } from './login/login.component';

@NgModule({
  declarations: [LoginComponent],
  imports: [
    CommonModule,
    OAuthModule.forRoot(),
    RouterModule.forChild([
      {
        path: 'login',
        component: LoginComponent
      }
    ])
  ],
  providers: [AuthenticationService, provideHttpClient(withInterceptorsFromDi())]
})
export class AuthenticationModule {
  static forRoot(): ModuleWithProviders<AuthenticationModule> {
    return {
      ngModule: AuthenticationModule,
      providers: [
        provideAppInitializer(() => {
          const initializerFn = authAppInitializerFactory(inject(AuthenticationService));
          return initializerFn();
        }),
        { provide: AuthConfig, useValue: authConfig },
        { provide: OAuthModuleConfig, useValue: authModuleConfig },
        { provide: OAuthStorage, useValue: localStorage }
      ]
    };
  }

  constructor(@Optional() @SkipSelf() parentModule: AuthenticationModule) {
    if (parentModule) {
      throw new Error('AuthenticationModule is already loaded. Import it in the AppModule only');
    }
  }
}

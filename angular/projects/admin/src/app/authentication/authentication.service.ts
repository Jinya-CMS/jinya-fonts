import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { OAuthErrorEvent, OAuthService } from 'angular-oauth2-oidc';
import { BehaviorSubject, combineLatest, Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';

@Injectable({ providedIn: 'root' })
export class AuthenticationService {
  private isAuthenticatedSubject$ = new BehaviorSubject<boolean>(false);
  public isAuthenticated$ = this.isAuthenticatedSubject$.asObservable();

  public isAuthenticated = false;

  private isDoneLoadingSubject$ = new BehaviorSubject<boolean>(false);
  public isDoneLoading$ = this.isDoneLoadingSubject$.asObservable();

  public canActivateProtectedRoutes$: Observable<boolean> = combineLatest([
    this.isAuthenticated$,
    this.isDoneLoading$
  ]).pipe(map((values) => values.every((b) => b)));

  private async navigateToLoginPage() {
    await this.router.navigateByUrl('login');
  }

  constructor(
    private oauthService: OAuthService,
    private router: Router
  ) {
    // Useful for debugging:
    this.oauthService.events.subscribe((event) => {
      if (event instanceof OAuthErrorEvent) {
        console.error('OAuthErrorEvent Object:', event);
      } else {
        console.warn('OAuthEvent Object:', event);
      }
    });

    this.isAuthenticated$.subscribe((value) => (this.isAuthenticated = value));

    this.oauthService.events.subscribe(() => {
      this.isAuthenticatedSubject$.next(this.oauthService.hasValidAccessToken());
    });
    this.isAuthenticatedSubject$.next(this.oauthService.hasValidAccessToken());

    this.oauthService.events
      .pipe(filter((e) => ['token_received'].includes(e.type)))
      .subscribe(() => this.oauthService.loadUserProfile());

    this.oauthService.events
      .pipe(filter((e) => ['session_terminated', 'session_error'].includes(e.type)))
      .subscribe(() => this.navigateToLoginPage());

    this.oauthService.setupAutomaticSilentRefresh();
  }

  public async runInitialLoginSequence() {
    try {
      await this.oauthService.loadDiscoveryDocument();
      await this.oauthService.tryLoginCodeFlow();
      if (this.oauthService.hasValidAccessToken()) {
        this.isDoneLoadingSubject$.next(true);

        if (this.router.url.includes('login')) {
          await this.router.navigateByUrl('font');
        }
      } else {
        await this.navigateToLoginPage();
      }
    } catch (e) {
      console.error(e);
      this.isDoneLoadingSubject$.next(true);
      await this.navigateToLoginPage();
    }
  }

  public login(targetUrl?: string) {
    this.oauthService.initLoginFlow(targetUrl || this.router.url);
  }

  public logout() {
    this.oauthService.logOut();
  }

  public refresh() {
    this.oauthService.silentRefresh();
  }

  public hasValidToken() {
    return this.oauthService.hasValidAccessToken();
  }

  public get accessToken() {
    return this.oauthService.getAccessToken();
  }

  public get refreshToken() {
    return this.oauthService.getRefreshToken();
  }

  public get identityClaims() {
    return this.oauthService.getIdentityClaims();
  }

  public get idToken() {
    return this.oauthService.getIdToken();
  }

  public get logoutUrl() {
    return this.oauthService.logoutUrl;
  }
}

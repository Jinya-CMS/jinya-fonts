import { Component } from '@angular/core';
import { Router, RouterLinkActive } from '@angular/router';
import { AuthenticationService } from './authentication/authentication.service';
import { ApiService } from './api/services/api.service';
import { FontSyncEvents, FontSyncService } from './fonts/font-sync.service';
import { Location } from '@angular/common';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
  standalone: false
})
export class AppComponent {
  protected readonly RouterLinkActive = RouterLinkActive;

  googleFontsSyncing = false;
  googleFontsSyncSuccess = false;
  googleFontsSyncFailure = false;

  constructor(
    protected authService: AuthenticationService,
    protected router: Router,
    protected location: Location,
    private apiService: ApiService,
    private fontSyncService: FontSyncService
  ) {
    fontSyncService.subscribe(FontSyncEvents.Start).subscribe(() => (this.googleFontsSyncing = true));
    fontSyncService.subscribe(FontSyncEvents.Success).subscribe(() => {
      this.googleFontsSyncSuccess = true;
      this.googleFontsSyncing = false;
    });
    fontSyncService.subscribe(FontSyncEvents.Error).subscribe(() => {
      this.googleFontsSyncFailure = true;
      this.googleFontsSyncing = false;
    });
  }
}

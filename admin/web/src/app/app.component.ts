import { Component } from '@angular/core';
import { Router, RouterLinkActive } from '@angular/router';
import { AuthenticationService } from './authentication/authentication.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  protected readonly RouterLinkActive = RouterLinkActive;

  constructor(
    protected authService: AuthenticationService,
    private router: Router
  ) {}

  goBack() {
    history.back();
  }
}

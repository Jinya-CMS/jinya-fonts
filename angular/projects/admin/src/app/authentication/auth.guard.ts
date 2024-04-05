import { CanActivateFn, Router } from '@angular/router';
import { inject } from '@angular/core';

import { AuthenticationService } from './authentication.service';

export const authGuard: CanActivateFn = () => {
  const auth = inject(AuthenticationService);
  const router = inject(Router);
  if (!auth.isAuthenticated) {
    return router.navigateByUrl('login');
  }

  return auth.isAuthenticated;
};

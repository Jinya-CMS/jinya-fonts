import { AuthenticationService } from './authentication.service';

export function authAppInitializerFactory(authService: AuthenticationService): () => Promise<void> {
  return () => authService.runInitialLoginSequence();
}

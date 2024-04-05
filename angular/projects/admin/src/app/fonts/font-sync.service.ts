import { Injectable } from '@angular/core';
import { ApiService } from '../api/services/api.service';
import { Observable, Subject } from 'rxjs';

export enum FontSyncEvents {
  Start,
  Success,
  Error
}

@Injectable({
  providedIn: 'root'
})
export class FontSyncService {
  constructor(private apiService: ApiService) {}

  private start = new Subject<void>();
  private success = new Subject<void>();
  private error = new Subject<void>();

  public subscribe(event: FontSyncEvents): Observable<void> {
    switch (event) {
      case FontSyncEvents.Start:
        return this.start.asObservable();
      case FontSyncEvents.Success:
        return this.success.asObservable();
      case FontSyncEvents.Error:
        return this.error.asObservable();
    }
  }

  public syncFonts() {
    this.start.next();
    this.apiService.syncGoogleFonts().subscribe({
      next: () => this.success.next(),
      error: () => this.error.next()
    });
  }
}

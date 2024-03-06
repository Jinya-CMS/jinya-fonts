import { NgModule } from '@angular/core';
import { AppComponent } from './app.component';
import { NgIf } from '@angular/common';
import { RouterLink, RouterLinkActive, RouterModule, RouterOutlet } from '@angular/router';
import { BrowserModule } from '@angular/platform-browser';
import { ApiModule } from './api/api.module';
import { UiModule } from './ui/ui.module';
import { HttpClientModule } from '@angular/common/http';
import { PagesModule } from './pages/pages.module';

@NgModule({
  declarations: [AppComponent],
  imports: [
    ApiModule.forRoot({ rootUrl: '' }),
    BrowserModule,
    NgIf,
    HttpClientModule,
    RouterLink,
    UiModule.forRoot(),
    PagesModule.forRoot(),
    RouterOutlet,
    RouterLinkActive,
    RouterModule.forRoot([], {
      bindToComponentInputs: true
    })
  ],
  bootstrap: [AppComponent]
})
export class AppModule {}

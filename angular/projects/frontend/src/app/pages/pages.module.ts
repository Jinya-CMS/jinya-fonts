import { ModuleWithProviders, NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { DetailPageComponent } from './detail-page/detail-page.component';
import { StartPageComponent } from './start-page/start-page.component';
import { ApiModule } from '../api/api.module';
import { UiModule } from '../ui/ui.module';
import { PreviewPanelComponent } from '../ui/preview-panel/preview-panel.component';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';

@NgModule({
  declarations: [StartPageComponent, DetailPageComponent],
  imports: [
    ApiModule.forRoot({ rootUrl: '' }),
    CommonModule,
    RouterModule.forChild([
      {
        path: 'font',
        component: DetailPageComponent
      },
      {
        path: '**',
        component: StartPageComponent
      }
    ]),
    UiModule,
    PreviewPanelComponent
  ],
  providers: [provideHttpClient(withInterceptorsFromDi())]
})
export class PagesModule {
  static forRoot(): ModuleWithProviders<PagesModule> {
    return {
      ngModule: PagesModule
    };
  }
}

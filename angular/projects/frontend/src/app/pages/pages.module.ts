import { ModuleWithProviders, NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { DetailPageComponent } from './detail-page/detail-page.component';
import { StartPageComponent } from './start-page/start-page.component';
import { ApiModule } from '../api/api.module';
import { HttpClientModule } from '@angular/common/http';
import { UiModule } from '../ui/ui.module';
import { PreviewPanelComponent } from '../ui/preview-panel/preview-panel.component';

@NgModule({
  declarations: [StartPageComponent, DetailPageComponent],
  imports: [
    ApiModule.forRoot({ rootUrl: '' }),
    HttpClientModule,
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
  ]
})
export class PagesModule {
  static forRoot(): ModuleWithProviders<PagesModule> {
    return {
      ngModule: PagesModule
    };
  }
}

import { ModuleWithProviders, NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { SearchbarComponent } from './searchbar/searchbar.component';
import { FilterPanelComponent } from './filter-panel/filter-panel.component';
import { FontCardComponent } from './font-card/font-card.component';
import { RouterLink } from '@angular/router';
import { PreviewPanelComponent } from './preview-panel/preview-panel.component';

@NgModule({
  declarations: [SearchbarComponent, FilterPanelComponent, FontCardComponent],
  exports: [FontCardComponent, FilterPanelComponent],
  imports: [CommonModule, ReactiveFormsModule, RouterLink, PreviewPanelComponent]
})
export class UiModule {
  static forRoot(): ModuleWithProviders<UiModule> {
    return {
      ngModule: UiModule
    };
  }
}

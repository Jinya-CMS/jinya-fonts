import { Component, OnInit } from '@angular/core';
import { ApiModule } from '../api/api.module';
import { ApiService } from '../api/services/api.service';
import { previewTexts, SettingsPanelComponent, WebfontFilter } from '../settings-panel/settings-panel.component';
import { HttpClientModule } from '@angular/common/http';
import { Webfont } from '../api/models/webfont';
import { FontCardComponent } from '../font-card/font-card.component';

@Component({
  selector: 'app-startpage',
  standalone: true,
  imports: [ApiModule, HttpClientModule, SettingsPanelComponent, FontCardComponent],
  templateUrl: './startpage.component.html',
  styleUrl: './startpage.component.scss'
})
export class StartpageComponent implements OnInit {
  constructor(private apiService: ApiService) {}

  webfonts: Webfont[] = [];
  filteredWebfonts: Webfont[] = [];
  loading = true;
  text = previewTexts.lorem;
  size = 24;

  ngOnInit(): void {
    this.apiService.apiFontGet().subscribe((value) => {
      this.webfonts = value;
      this.filteredWebfonts = value;
      this.loading = false;
    });
  }

  filter(filters: WebfontFilter) {
    this.filteredWebfonts = this.webfonts
      .filter((font) => font.name.toLowerCase().includes(filters.searchText))
      .filter((item) => {
        let result = false;
        if (filters.sansSerif) {
          result = result || item.category.toLowerCase() === 'sans serif';
        }
        if (filters.serif) {
          result = result || item.category.toLowerCase() === 'serif';
        }
        if (filters.handwriting) {
          result = result || item.category.toLowerCase() === 'handwriting';
        }
        if (filters.display) {
          result = result || item.category.toLowerCase() === 'display';
        }
        if (filters.monospace) {
          result = result || item.category.toLowerCase() === 'monospace';
        }

        return result;
      });
  }
}

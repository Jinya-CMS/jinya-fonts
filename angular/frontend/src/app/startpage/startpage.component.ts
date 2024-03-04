import { Component, OnInit } from '@angular/core';
import { ApiModule } from '../api/api.module';
import { ApiService } from '../api/services/api.service';
import { previewTexts, SettingsPanelComponent } from '../settings-panel/settings-panel.component';
import { HttpClientModule } from '@angular/common/http';
import { Webfont } from '../api/models/webfont';
import { FontCardComponent } from '../font-card/font-card.component';

interface WebfontFilter {
  searchText: string;
}

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
  filters: WebfontFilter = { searchText: '' };

  ngOnInit(): void {
    this.apiService.apiFontGet().subscribe((value) => {
      this.webfonts = value;
      this.filteredWebfonts = value;
      this.loading = false;
    });
  }

  search($event: string) {
    this.filters.searchText = $event;
    this.filter();
  }

  private filter() {
    this.filteredWebfonts = this.webfonts.filter((font) => font.name.toLowerCase().includes(this.filters.searchText));
  }
}

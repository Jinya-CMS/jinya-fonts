import { Component, HostListener, OnInit } from '@angular/core';
import { ApiModule } from '../api/api.module';
import { ApiService } from '../api/services/api.service';
import { previewTexts, SettingsPanelComponent, WebfontFilter } from '../settings-panel/settings-panel.component';
import { HttpClientModule } from '@angular/common/http';
import { Webfont } from '../api/models/webfont';
import { FontCardComponent } from '../font-card/font-card.component';

interface WebfontWithRow extends Webfont {
  rowStart: number;
  rowEnd: number;
  rowStartFull: number;
  rowEndFull: number;
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
  filteredWebfonts: WebfontWithRow[] = [];
  fontsToShow: WebfontWithRow[] = [];
  loading = true;
  text = previewTexts.lorem;
  size = 24;
  pageSize = Math.floor(window.innerHeight / (13 * 16));

  protected readonly Math = Math;

  ngOnInit(): void {
    this.apiService.apiFontGet().subscribe((value) => {
      this.webfonts = value;
      this.loading = false;
      this.filter({
        display: true,
        handwriting: true,
        monospace: true,
        sansSerif: true,
        searchText: '',
        serif: true
      });
      this.setFontsToShow();
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
      })
      .map((font, idx) => ({
        ...font,
        rowStart: Math.floor((idx === 0 ? 0 : idx / 4) + 1),
        rowEnd: Math.floor((idx === 0 ? 0 : idx / 4) + 2),
        rowStartFull: idx + 1,
        rowEndFull: idx + 2
      }));
    this.setFontsToShow();
    window.scrollTo({
      top: 0
    });
  }

  setFontsToShow() {
    let page = Math.floor((window.scrollY + window.innerHeight) / window.innerHeight) - 1;
    if (page < 0) {
      page = 0;
    }

    const startIndex = page * this.pageSize;
    const endIndex = startIndex * 4 + this.pageSize * 4 * 2 * 2;
    this.fontsToShow = this.filteredWebfonts.slice(startIndex * 4, endIndex);
  }

  @HostListener('window:scroll')
  windowScroll() {
    this.setFontsToShow();
  }
}
